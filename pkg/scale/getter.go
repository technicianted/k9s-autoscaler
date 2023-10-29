// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package scale

import (
	"k9s-autoscaler/pkg/scale/types"

	"k8s.io/client-go/scale"
)

var (
	_ scale.ScalesGetter = &getter{}
)

// Simple adapter for k8s ScalingClient.
type getter struct {
	scaler types.ScalingClient
}

func NewGetter(scaler types.ScalingClient) scale.ScalesGetter {
	return &getter{
		scaler: scaler,
	}
}

func (g *getter) Scales(namespace string) scale.ScaleInterface {
	return NewScaler(namespace, g.scaler)
}
