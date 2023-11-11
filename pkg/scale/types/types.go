// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"

	prototypes "k9s-autoscaler/pkg/proto"
)

//go:generate mockgen -package mocks -destination ../mocks/scale.go -source $GOFILE

// Scaling client implements a scaling adapter for kubernetes hpa.
type ScalingClient interface {
	// Scales a scaleTarget to given target for autoscaler name and namespace.
	SetScaleTarget(ctx context.Context, name, namespace string, scaleTarget *prototypes.AutoscalerTarget, target *prototypes.ScaleSpec) error
	// Gets current scale for scaleTarget for autoscaler name and namespace.
	GetScale(ctx context.Context, name, namespace string, scaleTarget *prototypes.AutoscalerTarget) (*prototypes.Scale, error)
}
