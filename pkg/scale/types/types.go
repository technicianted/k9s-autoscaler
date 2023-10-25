// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"

	prototypes "k9s-autoscaler/pkg/proto"
)

//go:generate mockgen -package mocks -destination ../mocks/scale.go -source $GOFILE

type ScalingClient interface {
	SetScaleTarget(ctx context.Context, name, namespace string, target *prototypes.ScaleSpec) error
	GetScale(ctx context.Context, name, namespace string) (*prototypes.Scale, error)
}
