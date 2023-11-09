// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"
	"time"
)

//go:generate mockgen -package mocks -destination ../mocks/metrics.go -source $GOFILE

// An interface defining a provider adapter implementation that can fetch metric
// values.
type MetricsClient interface {
	// Get metric values metricName in namespace using selector keys and values. Returns an array
	// of values and values timestamp.
	GetMetric(ctx context.Context, metricName string, namespace string, selector map[string]string) ([]int64, time.Time, error)
}
