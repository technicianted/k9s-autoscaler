// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"k9s-autoscaler/pkg/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricNameLabel = "metric"
)

var (
	metricLatencyMetric = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: common.MetricsNamespace,
			Subsystem: "metrics",
			Name:      "get_latency",
			Help:      "Time to get value of a metric.",
			Buckets:   prometheus.ExponentialBucketsRange(0.001, 60.0, 16),
		},
		[]string{common.MetricsAutoscalerNamespaceLabel, metricNameLabel, common.MetricsErrorLabel})
)
