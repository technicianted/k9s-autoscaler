// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package scale

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/scale/types"
	storagetypes "k9s-autoscaler/pkg/storage/types"

	autoscalingapi "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/scale"
	"k8s.io/klog/v2"
)

// An adapter for k8s ScalingClient interface.
type scaler struct {
	scale.ScaleInterface

	namespace        string
	autoscalerGetter storagetypes.AutoscalerGetter
	scaler           types.ScalingClient
}

// Create a new scaler adapter for namespace that uses s to pass scaling calls.
func NewScaler(namespace string, autoscalerGetter storagetypes.AutoscalerGetter, s types.ScalingClient) scale.ScaleInterface {
	return &scaler{
		namespace:        namespace,
		autoscalerGetter: autoscalerGetter,
		scaler:           s,
	}
}

// Get fetches the scale of the given scalable resource.
func (s *scaler) Get(ctx context.Context, resource schema.GroupResource, name string, opts metav1.GetOptions) (*autoscalingapi.Scale, error) {
	klog.V(1).InfoS("get scale", "name", name)

	as, err := s.autoscalerGetter.Get(name, s.namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get autoscaler: %v", err)
	}

	opTimer := time.Now()
	scale, err := s.scaler.GetScale(ctx, name, s.namespace, as.Spec.Target)
	if err != nil {
		scaleLatencyMetric.WithLabelValues(name, s.namespace, opGet, "true").Observe(time.Since(opTimer).Seconds())
		return nil, err
	}

	klog.V(1).InfoS("get scale", "name", name, "scale", scale)

	var errs field.ErrorList
	if scale.Spec == nil {
		errs = append(errs, &field.Error{
			Type:   field.ErrorTypeRequired,
			Field:  "spec",
			Detail: "scale spec",
		})
	}
	if scale.Status == nil {
		errs = append(errs, &field.Error{
			Type:   field.ErrorTypeRequired,
			Field:  "status",
			Detail: "scale status",
		})
	}
	if len(errs) > 0 {
		scaleLatencyMetric.WithLabelValues(name, s.namespace, opGet, "true").Observe(time.Since(opTimer).Seconds())
		return nil, errors.NewInvalid(schema.GroupKind{
			Group: resource.Group,
			Kind:  "scale",
		},
			name,
			errs)
	}

	scaleLatencyMetric.WithLabelValues(name, s.namespace, opGet, "").Observe(time.Since(opTimer).Seconds())

	return &autoscalingapi.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.namespace,
		},
		Spec: autoscalingapi.ScaleSpec{
			Replicas: int32(scale.Spec.Desired),
		},
		Status: autoscalingapi.ScaleStatus{
			Replicas: int32(scale.Status.Current),
			Selector: EncodePodLabels(scale.Status.Current),
		},
	}, nil
}

// Update updates the scale of the given scalable resource.
func (s *scaler) Update(ctx context.Context, resource schema.GroupResource, scale *autoscalingapi.Scale, opts metav1.UpdateOptions) (*autoscalingapi.Scale, error) {
	klog.V(1).InfoS("update scale", "name", scale.Name, "target", scale.Spec.Replicas)

	as, err := s.autoscalerGetter.Get(scale.Name, s.namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get autoscaler: %v", err)
	}

	opTimer := time.Now()

	err = s.scaler.SetScaleTarget(
		ctx,
		scale.Name,
		s.namespace,
		as.Spec.Target,
		&prototypes.ScaleSpec{
			Desired: scale.Spec.Replicas,
		})
	if err != nil {
		scaleLatencyMetric.WithLabelValues(scale.Name, s.namespace, opSet, "true").Observe(time.Since(opTimer).Seconds())
		return nil, err
	}
	scaleLatencyMetric.WithLabelValues(scale.Name, s.namespace, opSet, "").Observe(time.Since(opTimer).Seconds())

	return scale, nil
}

func EncodePodLabels(replicas int32) string {
	return fmt.Sprintf("replicas=%d", replicas)
}

func DecodePodLabels(selector labels.Selector) int32 {
	l := selector.String()
	parts := strings.Split(l, "=")
	if len(parts) != 2 {
		panic(fmt.Sprintf("unexpected scale pod selector format: %v", selector))
	}
	if parts[0] != "replicas" {
		panic(fmt.Sprintf("unexpected scale pod selector label: %v", selector))
	}
	replicas, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Sprintf("unexpected scale pod selector replicas: %v", selector))
	}

	return int32(replicas)
}
