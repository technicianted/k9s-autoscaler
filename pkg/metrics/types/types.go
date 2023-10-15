// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"
	"time"
)

//go:generate mockgen -package mocks -destination ../mocks/metrics.go -source $GOFILE

type MetricsClient interface {
	// GetExternalMetric gets all the values of a given external metric
	// that match the specified selector.
	GetMetric(ctx context.Context, autoscaler, namespace, metricName string) ([]int64, time.Time, error)
}
