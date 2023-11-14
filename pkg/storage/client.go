// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/scale"
	"k9s-autoscaler/pkg/storage/metrics"
	"k9s-autoscaler/pkg/storage/types"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2"
	"k8s.io/klog/v2"
)

type autoscalerEntry struct {
	autoscaler *prototypes.Autoscaler
	hpa        *v2.HorizontalPodAutoscaler
}

// An autoscaler client that implements storage mapping between K9s and K8s
// autoscalers by providing adapters to Kubernetes lister, informer, etc.
type Client struct {
	sync.RWMutex
	apiv2.HorizontalPodAutoscalersGetter

	statusUpdateHandler       types.AutoscalerStatusUpdateHandler
	autoscalerByNamespaceName map[string]map[string]*autoscalerEntry
	watchersByNamespace       map[string]types.AutoscalerStatusUpdateHandler
	watchesByNamespace        map[string]map[*autoscalerWatch]bool
	watchesByWatchNamespace   map[*autoscalerWatch]string
}

// Create a new client that uses statusUpdateHandler to propagate changes in
// underlying autoscaler status by the HPA.
func NewClient(statusUpdatedHandler types.AutoscalerStatusUpdateHandler) (*Client, error) {
	c := &Client{
		statusUpdateHandler:       statusUpdatedHandler,
		autoscalerByNamespaceName: make(map[string]map[string]*autoscalerEntry),
		watchersByNamespace:       make(map[string]types.AutoscalerStatusUpdateHandler),
		watchesByNamespace:        make(map[string]map[*autoscalerWatch]bool),
		watchesByWatchNamespace:   make(map[*autoscalerWatch]string),
	}
	if err := metrics.RegisterMetricsCollector(c); err != nil {
		return nil, err
	}

	return c, nil
}

// Implement HorizontalPodAutoscalersGetter k8s interface.
func (c *Client) HorizontalPodAutoscalers(namespace string) apiv2.HorizontalPodAutoscalerInterface {
	return newNamespacedClient(c, namespace)
}

// Returns a list of all available autoscalers. Returned objects are clones
// of the internal storage so they are safe to be modified without affecting
// internal state.
// Implements AutoscalerCRUDder.
func (c *Client) List() ([]*prototypes.Autoscaler, error) {
	c.RLock()
	defer c.RUnlock()

	autoscalers := make([]*prototypes.Autoscaler, 0)
	for _, autoscalersByName := range c.autoscalerByNamespaceName {
		for _, autoscalerEntry := range autoscalersByName {
			copy := proto.Clone(autoscalerEntry.autoscaler).(*prototypes.Autoscaler)
			autoscalers = append(autoscalers, copy)
		}
	}

	return autoscalers, nil
}

// Get an autoscaler by name and namespace.
func (c *Client) Get(name, namespace string) (*prototypes.Autoscaler, error) {
	c.RLock()
	defer c.RUnlock()

	if autoscalersByName, ok := c.autoscalerByNamespaceName[namespace]; ok {
		if autoscaler, ok := autoscalersByName[name]; ok {
			return autoscaler.autoscaler, nil
		}
	}

	return nil, fmt.Errorf("autoscaler %s namespace %s not found", name, namespace)
}

// Adds a new autoscaler. All defined k8s watches will be notified.
// Implements AutoscalerCRUDder.
func (c *Client) Add(autoscaler *prototypes.Autoscaler) error {
	klog.V(0).InfoS("adding new autoscaler", "autoscaler", autoscaler)

	c.Lock()
	defer c.Unlock()

	if _, ok := c.autoscalerByNamespaceName[autoscaler.Namespace]; !ok {
		c.autoscalerByNamespaceName[autoscaler.Namespace] = make(map[string]*autoscalerEntry)
	}
	if _, ok := c.autoscalerByNamespaceName[autoscaler.Namespace][autoscaler.Name]; ok {
		return fmt.Errorf("already exists: %s", autoscaler.Name)
	}

	hpa, err := autoscalerToHPA(autoscaler)
	if err != nil {
		return err
	}
	entry := &autoscalerEntry{
		autoscaler: autoscaler,
		hpa:        hpa,
	}
	c.autoscalerByNamespaceName[autoscaler.Namespace][autoscaler.Name] = entry

	c.updateWatchesAddedLocked(entry)

	return nil
}

// Updates an existing autoscaler. All defined k8s watches will be notified.
// Implements AutoscalerCRUDder.
func (c *Client) Update(autoscaler *prototypes.Autoscaler) error {
	klog.V(0).InfoS("updating autoscaler", "autoscaler", autoscaler)

	c.Lock()
	defer c.Unlock()

	if _, ok := c.autoscalerByNamespaceName[autoscaler.Namespace]; !ok {
		return errors.NewNotFound(v2.Resource("horizontalpodautoscaler"), autoscaler.Name)
	}
	if _, ok := c.autoscalerByNamespaceName[autoscaler.Namespace][autoscaler.Name]; !ok {
		return errors.NewNotFound(v2.Resource("horizontalpodautoscaler"), autoscaler.Name)
	}

	hpa, err := autoscalerToHPA(autoscaler)
	if err != nil {
		return err
	}
	entry := &autoscalerEntry{
		autoscaler: autoscaler,
		hpa:        hpa,
	}
	c.autoscalerByNamespaceName[autoscaler.Namespace][autoscaler.Name] = entry

	c.updateWatchesModifiedLocked(entry)

	return nil
}

// Deletes an existing autoscaler with name and namespace. All defined k8s
// watches will be notified.
// Implements AutoscalerCRUDder.
func (c *Client) Delete(name, namespace string) error {
	klog.V(0).InfoS("deleting autoscaler", "name", name, "namespace", namespace)
	c.Lock()
	defer c.Unlock()

	if _, ok := c.autoscalerByNamespaceName[namespace]; !ok {
		return errors.NewNotFound(v2.Resource("horizontalpodautoscaler"), name)
	}
	if entry, ok := c.autoscalerByNamespaceName[namespace][name]; !ok {
		return errors.NewNotFound(v2.Resource("horizontalpodautoscaler"), name)
	} else {
		delete(c.autoscalerByNamespaceName[namespace], name)

		c.updateWatchesDeletedLocked(entry)
	}

	return nil
}

func (c *Client) updateWatchesAddedLocked(entry *autoscalerEntry) {
	c.updatedWatchesLocked(entry.autoscaler.Namespace, func(w *autoscalerWatch) {
		w.add(entry.hpa)
	})
}

func (c *Client) updateWatchesModifiedLocked(entry *autoscalerEntry) {
	c.updatedWatchesLocked(entry.autoscaler.Namespace, func(w *autoscalerWatch) {
		w.update(entry.hpa)
	})
}

func (c *Client) updateWatchesDeletedLocked(entry *autoscalerEntry) {
	c.updatedWatchesLocked(entry.autoscaler.Namespace, func(w *autoscalerWatch) {
		w.delete(entry.hpa)
	})
}

func (c *Client) updatedWatchesLocked(namespace string, watchFunc func(*autoscalerWatch)) {
	if namespace != v1.NamespaceAll {
		if watches, ok := c.watchesByNamespace[namespace]; ok {
			for watch := range watches {
				watchFunc(watch)
			}
		}
	}
	if watches, ok := c.watchesByNamespace[v1.NamespaceAll]; ok {
		for watch := range watches {
			watchFunc(watch)
		}
	}
}

func autoscalerToHPA(autoscaler *prototypes.Autoscaler) (*v2.HorizontalPodAutoscaler, error) {
	if len(autoscaler.Name) == 0 {
		return nil, fmt.Errorf("name is required")
	}
	if len(autoscaler.Namespace) == 0 {
		return nil, fmt.Errorf("namespace is rquired")
	}
	if autoscaler.Spec == nil || len(autoscaler.Spec.Metrics) == 0 {
		return nil, fmt.Errorf("no metrics")
	}

	metrics := make([]v2.MetricSpec, len(autoscaler.Spec.Metrics))
	for i := 0; i < len(autoscaler.Spec.Metrics); i++ {
		metric := autoscaler.Spec.Metrics[i]
		metrics[i] = v2.MetricSpec{
			Type: v2.ExternalMetricSourceType,
			External: &v2.ExternalMetricSource{
				Metric: v2.MetricIdentifier{
					Name:     metric.Name,
					Selector: EncodeMetricHPA(autoscaler.Name),
				},
				Target: v2.MetricTarget{
					Type:  v2.ValueMetricType,
					Value: resource.NewQuantity(metric.Target, resource.DecimalSI),
				},
			},
		}
	}
	behavior, err := autoscalerBehaviorToHPA(autoscaler.Spec.Behavior)
	if err != nil {
		return nil, err
	}
	return &v2.HorizontalPodAutoscaler{
		TypeMeta: v1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:            autoscaler.Name,
			Namespace:       autoscaler.Namespace,
			ResourceVersion: autoscaler.Version,
		},
		Spec: v2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2.CrossVersionObjectReference{
				Kind:       scale.ScalingResourceKind,
				APIVersion: filepath.Join(scale.ScalingResourceGroup, scale.ScalingResourceVersion),
				Name:       autoscaler.Name,
			},
			MinReplicas: int32ToInt32Pointer(autoscaler.Spec.Min),
			MaxReplicas: int32(autoscaler.Spec.Max),
			Metrics:     metrics,
			Behavior:    behavior,
		},
		Status: v2.HorizontalPodAutoscalerStatus{
			LastScaleTime: &v1.Time{},
		},
	}, nil
}

func autoscalerBehaviorToHPA(behavior *prototypes.Behavior) (*v2.HorizontalPodAutoscalerBehavior, error) {
	if behavior == nil {
		return nil, nil
	}
	scaleUp, err := autoscalerScalingRulesToHPA(behavior.ScaleUp)
	if err != nil {
		return nil, err
	}
	scaleDown, err := autoscalerScalingRulesToHPA(behavior.ScaleDown)
	if err != nil {
		return nil, err
	}
	hpaBehavior := v2.HorizontalPodAutoscalerBehavior{
		ScaleUp:   scaleUp,
		ScaleDown: scaleDown,
	}

	return &hpaBehavior, nil
}

func autoscalerScalingRulesToHPA(rules *prototypes.ScalingRules) (*v2.HPAScalingRules, error) {
	if rules == nil {
		return nil, nil
	}

	hpaRules := v2.HPAScalingRules{
		StabilizationWindowSeconds: rules.StabilizationWindowSeconds,
	}

	var policy v2.ScalingPolicySelect
	switch rules.SelectPolicy {
	case prototypes.ScalingRules_Max:
		policy = v2.MaxChangePolicySelect
	case prototypes.ScalingRules_Min:
		policy = v2.MinChangePolicySelect
	case prototypes.ScalingRules_Disabled:
		policy = v2.DisabledPolicySelect
	default:
		return nil, fmt.Errorf("unexpected scaling policy: %v", rules.SelectPolicy)
	}
	hpaRules.SelectPolicy = &policy

	for _, rule := range rules.Policies {
		valueType := v2.PercentScalingPolicy
		if rule.ValueType == prototypes.ScalingPolicy_Units {
			valueType = v2.PodsScalingPolicy
		}
		hpaRules.Policies = append(hpaRules.Policies, v2.HPAScalingPolicy{
			Type:          valueType,
			Value:         rule.Value,
			PeriodSeconds: rule.PeriodSeconds,
		})
	}

	return &hpaRules, nil
}

func hpaStatusToAutoScaler(status v2.HorizontalPodAutoscalerStatus) (*prototypes.AutoscalerStatus, error) {
	autoscalerStatus := prototypes.AutoscalerStatus{
		CurrentScale: int32ToInt32Pointer(status.CurrentReplicas),
		DesiredScale: status.DesiredReplicas,
	}
	for _, condition := range status.Conditions {
		var typ prototypes.Condition_ConditionType
		switch condition.Type {
		case v2.ScalingActive:
			typ = prototypes.Condition_ScalingActive
		case v2.ScalingLimited:
			typ = prototypes.Condition_ScalingLimited
		case v2.AbleToScale:
			typ = prototypes.Condition_AbleToScale
		default:
			return nil, fmt.Errorf("unexpected condition type: %v", condition.Type)
		}
		autoscalerStatus.Conditions = append(autoscalerStatus.Conditions, &prototypes.Condition{
			Type:               typ,
			Status:             string(condition.Status),
			LastTransitionTime: timestamppb.New(condition.LastTransitionTime.Time),
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}

	if status.LastScaleTime != nil {
		autoscalerStatus.LastScaleTime = timestamppb.New(status.LastScaleTime.Time)
	} else {
		autoscalerStatus.LastScaleTime = timestamppb.New(time.Time{})
	}

	return &autoscalerStatus, nil
}
