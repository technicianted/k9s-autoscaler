// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"

	prototypes "k9s-autoscaler/pkg/proto"
)

//go:generate mockgen -package mocks -destination ../mocks/events.go -source $GOFILE

type EventCreator interface {
	Create(ctx context.Context, name, namespace string, event *prototypes.AutoscalerEvent) error
}
