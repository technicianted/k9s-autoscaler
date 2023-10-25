// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import prototypes "k9s-autoscaler/pkg/proto"

//go:generate mockgen -package mocks -destination ../mocks/storage.go -source $GOFILE

const (
	ScaleResourceAPIVersion = ""
	ScaleResrouceKind       = ""
)

type AutoscalerCRUDder interface {
	List() ([]*prototypes.Autoscaler, error)
	Add(autoscaler *prototypes.Autoscaler) error
	Update(autoscaler *prototypes.Autoscaler) error
	Delete(name, namespace string) error
}

type AutoscalerStatusUpdateHandler interface {
	AutoscalerStatusUpdated(autoscaler *prototypes.Autoscaler)
}
