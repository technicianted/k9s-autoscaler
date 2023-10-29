// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"

	prototypes "k9s-autoscaler/pkg/proto"
)

//go:generate mockgen -package mocks -destination ../mocks/events.go -source $GOFILE

// Defines a provider adapter interface that will be called whenever the underlying
// HPA emits an event via k8s recorder interface.
type EventCreator interface {
	// Create a new event for autoscaler name and namespace.
	Create(ctx context.Context, name, namespace string, event *prototypes.AutoscalerEvent) error
}
