// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import prototypes "k9s-autoscaler/pkg/proto"

//go:generate mockgen -package mocks -destination ../mocks/storage.go -source $GOFILE

const (
	// A dummy API version for our scaling resource.
	ScaleResourceAPIVersion = ""
	ScaleResrouceKind       = ""
)

type AutoscalerGetter interface {
	// Get a list of defiend autosclaers.
	List() ([]*prototypes.Autoscaler, error)
	// Get an autoscaler by name and namespace.
	Get(name, namespace string) (*prototypes.Autoscaler, error)
}

// An interface to abstract CRUD operations for autoscalers. It can be used
// by storage providers to adapt to external APIs.
// All calls are expected to return from local cache and are nonblocking.
type AutoscalerCRUDder interface {
	AutoscalerGetter
	// Add new new autoscaler. Returns an error of already exists.
	Add(autoscaler *prototypes.Autoscaler) error
	// Updates an autoscaler. Returns an error if not exists.
	Update(autoscaler *prototypes.Autoscaler) error
	// Deletes an autoscaler. Return an error if not exists.
	Delete(name, namespace string) error
}

// A callback interface that defines a way for external adapters to receive
// autoscaler status update.
type AutoscalerStatusUpdateHandler interface {
	// Autoscaler status has been updated.
	AutoscalerStatusUpdated(autoscaler *prototypes.Autoscaler)
}
