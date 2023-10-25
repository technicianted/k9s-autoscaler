// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"context"
	"fmt"

	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	apiv2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2"
	"k8s.io/klog/v2"
)

type namespacedClient struct {
	apiv2.HorizontalPodAutoscalerInterface

	client    *Client
	namespace string
}

func newNamespacedClient(client *Client, namespace string) *namespacedClient {
	return &namespacedClient{
		client:    client,
		namespace: namespace,
	}
}

func (c *namespacedClient) UpdateStatus(ctx context.Context, horizontalPodAutoscaler *v2.HorizontalPodAutoscaler, opts v1.UpdateOptions) (*v2.HorizontalPodAutoscaler, error) {
	if len(opts.DryRun) > 0 {
		return nil, errors.NewForbidden(v2.Resource("horizontalpodautoscaler"), horizontalPodAutoscaler.Name, fmt.Errorf("dryrun is not allowed"))
	}
	if len(opts.FieldValidation) > 0 {
		return nil, errors.NewForbidden(v2.Resource("horizontalpodautoscaler"), horizontalPodAutoscaler.Name, fmt.Errorf("field validation is not allowed"))
	}
	if c.namespace != horizontalPodAutoscaler.Namespace {
		return nil, errors.NewConflict(
			v2.Resource("horizontalpodautoscaler"),
			horizontalPodAutoscaler.Name,
			fmt.Errorf("client namespace mismatch: %s!=%s", horizontalPodAutoscaler.Namespace, c.namespace))
	}

	klog.V(1).InfoS("updating status", "hpa", horizontalPodAutoscaler.ObjectMeta, "status", horizontalPodAutoscaler.Status)
	c.client.Lock()
	defer c.client.Unlock()

	entry, err := c.getEntryLocked(horizontalPodAutoscaler.Name)
	if err != nil {
		return nil, err
	}

	status, err := hpaStatusToAutoScaler(horizontalPodAutoscaler.Status)
	if err != nil {
		return nil, err
	}
	entry.autoscaler.Status = status
	horizontalPodAutoscaler.Status.DeepCopyInto(&entry.hpa.Status)

	if watcher, ok := c.client.watchersByNamespace[c.namespace]; ok {
		go func() {
			watcher.AutoscalerStatusUpdated(entry.autoscaler)
		}()
	}

	c.client.statusUpdateHandler.AutoscalerStatusUpdated(entry.autoscaler)

	c.client.updateWatchesModifiedLocked(entry)

	return entry.hpa, nil
}

func (c *namespacedClient) Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.HorizontalPodAutoscaler, error) {
	c.client.RLock()
	defer c.client.RUnlock()

	entry, err := c.getEntryLocked(name)
	if err != nil {
		return nil, err
	}

	return entry.hpa, nil
}

func (c *namespacedClient) List(ctx context.Context, opts v1.ListOptions) (*v2.HorizontalPodAutoscalerList, error) {
	klog.V(0).InfoS("listing hpas", "namespace", c.namespace)

	c.client.RLock()
	defer c.client.RUnlock()

	list := &v2.HorizontalPodAutoscalerList{}
	for namespace, autoscalersByName := range c.client.autoscalerByNamespaceName {
		if c.namespace == v1.NamespaceAll || namespace == c.namespace {
			for _, entry := range autoscalersByName {
				list.Items = append(list.Items, *entry.hpa)
			}
		}
	}

	klog.V(1).InfoS("list", "list", list)
	return list, nil
}

func (c *namespacedClient) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	klog.V(0).InfoS("creating new watch", "namespace", c.namespace)

	c.client.Lock()
	defer c.client.Unlock()

	watch := newWatch(c)
	if _, ok := c.client.watchesByNamespace[c.namespace]; !ok {
		c.client.watchesByNamespace[c.namespace] = make(map[*autoscalerWatch]bool)
	}
	c.client.watchesByNamespace[c.namespace][watch] = true

	return watch, nil
}

func (c *namespacedClient) removeWatch(watch *autoscalerWatch) {
	klog.V(0).InfoS("removing watch", "namespace", c.namespace)

	c.client.Lock()
	defer c.client.Unlock()

	if watches, ok := c.client.watchesByNamespace[c.namespace]; !ok {
		klog.V(0).InfoS("watch namespace not found", "namespace", c.namespace)
		return
	} else if _, ok := watches[watch]; !ok {
		klog.V(0).InfoS("watch not found", "namespace", c.namespace)
		return
	} else {
		delete(watches, watch)
	}
}

func (c *namespacedClient) getEntryLocked(name string) (*autoscalerEntry, error) {
	if _, ok := c.client.autoscalerByNamespaceName[c.namespace]; !ok {
		return nil, errors.NewNotFound(v2.Resource("horizontalpodautoscaler"), name)
	}
	if entry, ok := c.client.autoscalerByNamespaceName[c.namespace][name]; !ok {
		return nil, errors.NewNotFound(v2.Resource("horizontalpodautoscaler"), name)
	} else {
		return entry, nil
	}
}
