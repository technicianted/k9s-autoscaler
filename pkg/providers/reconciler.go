// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package providers

import (
	"errors"
	"fmt"

	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers/types"
	storagetypes "k9s-autoscaler/pkg/storage/types"

	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

type reconciler struct {
	client storagetypes.AutoscalerCRUDder
}

// Creates a new reconciler that uses client to reflect updates.
func NewReconciler(client storagetypes.AutoscalerCRUDder) types.Reconciler {
	return &reconciler{
		client: client,
	}
}

func (r *reconciler) Reconcile(autoscalers []*prototypes.Autoscaler) error {
	type autoscalerKey struct {
		Name      string
		Namespace string
	}

	list, err := r.client.List()
	if err != nil {
		return fmt.Errorf("failed to fetch autoscalers list: %v", err)
	}

	existing := make(map[autoscalerKey]*prototypes.Autoscaler)
	for _, autoscaler := range list {
		existing[autoscalerKey{autoscaler.Name, autoscaler.Namespace}] = autoscaler
	}
	incoming := make(map[autoscalerKey]*prototypes.Autoscaler)
	for _, autoscaler := range autoscalers {
		incoming[autoscalerKey{autoscaler.Name, autoscaler.Namespace}] = autoscaler
	}

	deleted := make(map[autoscalerKey]*prototypes.Autoscaler)
	updated := make(map[autoscalerKey]*prototypes.Autoscaler)
	added := make(map[autoscalerKey]*prototypes.Autoscaler)

	// deleted
	for id, existingAutoscaler := range existing {
		if autoscaler, ok := incoming[id]; !ok {
			deleted[id] = autoscaler
		} else if !AutoscalerEqual(existingAutoscaler, autoscaler) {
			updated[id] = autoscaler
		}
	}
	// add
	for id, autoscaler := range incoming {
		if _, ok := existing[id]; !ok {
			added[id] = autoscaler
		}
	}

	var errs []error
	for id := range deleted {
		if err := r.client.Delete(id.Name, id.Namespace); err != nil {
			klog.InfoS("reconciler failed to delete autoscaler", "id", id, "error", err)
			errs = append(errs, err)
		}
	}
	for id, autoscaler := range updated {
		if err := r.client.Update(autoscaler); err != nil {
			klog.InfoS("reconciler failed to update autoscaler", "id", id, "error", err)
			errs = append(errs, err)
		}
	}
	for id, autoscaler := range added {
		if err := r.client.Add(autoscaler); err != nil {
			klog.InfoS("reconciler failed to add autoscaler", "id", id, "error", err)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func AutoscalerEqual(a1, a2 *prototypes.Autoscaler) bool {
	if a1.Name != a2.Name || a1.Namespace != a2.Namespace || a1.Version != a2.Version {
		return false
	}

	return proto.Equal(a1.Spec, a2.Spec)
}
