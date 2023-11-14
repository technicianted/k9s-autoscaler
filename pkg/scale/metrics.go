package scale

import (
	"k9s-autoscaler/pkg/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	opLabel = "op"
	opGet   = "get"
	opSet   = "set"
)

var (
	scaleLatencyMetric = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: common.MetricsNamespace,
			Subsystem: "scale",
			Name:      "op_latency",
			Help:      "Time to make a scale operation.",
			Buckets:   prometheus.ExponentialBucketsRange(0.001, 60.0, 16),
		},
		[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel, opLabel, common.MetricsErrorLabel})
)
