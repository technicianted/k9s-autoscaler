// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
)

//go:generate mockgen -package mocks -destination ../mocks/metrics.go -source $GOFILE

// An interface defining a provider adapter implementation that can fetch metric
// values.
type MetricsClient interface {
	// Get metric values metricName in namespace using opaque provider configs. Returns an array
	// of values and values timestamp.
	GetMetric(ctx context.Context, metricName, autoscalerName, namespace string, config *anypb.Any) ([]int64, time.Time, error)
}
