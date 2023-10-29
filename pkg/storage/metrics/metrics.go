// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"k9s-autoscaler/pkg/common"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2"
	"k8s.io/klog/v2"
)

const (
	metricNameLabel = "metric"
)

var (
	UpdateStatusLatencyMetric = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: common.MetricsNamespace,
			Subsystem: "status",
			Name:      "update_latency",
			Help:      "Time to update status of an autoscaler.",
			Buckets:   prometheus.ExponentialBucketsRange(0.001, 5.0, 10),
		},
		[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel, common.MetricsErrorLabel})
)

var (
	_ prometheus.Collector = &metricsCollector{}
)

type metricsCollector struct {
	sync.Mutex

	getter autoscalingv2.HorizontalPodAutoscalersGetter

	specMinMetric           *prometheus.GaugeVec
	specMaxMetric           *prometheus.GaugeVec
	specMetricTarget        *prometheus.GaugeVec
	scaleStateDesiredMetric *prometheus.GaugeVec
	scaleStateCurrentMetric *prometheus.GaugeVec
	metricsCurrentMetric    *prometheus.GaugeVec
}

func RegisterMetricsCollector(getter autoscalingv2.HorizontalPodAutoscalersGetter) error {
	collector := metricsCollector{
		getter: getter,

		specMinMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: common.MetricsNamespace,
				Subsystem: "spec",
				Name:      "min_scale",
				Help:      "Minimum scale spec.",
			},
			[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel}),
		specMaxMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: common.MetricsNamespace,
				Subsystem: "spec",
				Name:      "max_scale",
				Help:      "Maximum scale spec.",
			},
			[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel}),
		specMetricTarget: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: common.MetricsNamespace,
				Subsystem: "spec",
				Name:      "metric_target",
				Help:      "Metric scale target.",
			},
			[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel, metricNameLabel}),
		scaleStateDesiredMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: common.MetricsNamespace,
				Subsystem: "status",
				Name:      "desired_scale",
				Help:      "Desired scale.",
			},
			[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel}),
		scaleStateCurrentMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: common.MetricsNamespace,
				Subsystem: "status",
				Name:      "current_scale",
				Help:      "Current scale.",
			},
			[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel}),
		metricsCurrentMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: common.MetricsNamespace,
				Subsystem: "status",
				Name:      "metric_value",
				Help:      "Current metric value.",
			},
			[]string{common.MetricsAutoscaleNameLabel, common.MetricsAutoscalerNamespaceLabel, metricNameLabel}),
	}

	return prometheus.Register(&collector)
}

func (c *metricsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.Lock()
	defer c.Unlock()

	c.specMinMetric.Describe(ch)
	c.specMaxMetric.Describe(ch)
	c.specMetricTarget.Describe(ch)
	c.scaleStateDesiredMetric.Describe(ch)
	c.scaleStateCurrentMetric.Describe(ch)
	c.metricsCurrentMetric.Describe(ch)
}

func (c *metricsCollector) Collect(ch chan<- prometheus.Metric) {
	c.Lock()
	defer c.Unlock()

	c.specMinMetric.Reset()
	c.specMaxMetric.Reset()
	c.specMetricTarget.Reset()
	c.scaleStateDesiredMetric.Reset()
	c.scaleStateCurrentMetric.Reset()
	c.metricsCurrentMetric.Reset()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	list, err := c.getter.HorizontalPodAutoscalers(v1.NamespaceAll).List(ctx, apimetav1.ListOptions{})
	if err != nil {
		klog.InfoS("storage client metrics collector failed to list autoscalers", "error", err)
		return
	}

	for _, autoscaler := range list.Items {
		if autoscaler.Spec.MinReplicas != nil {
			c.specMinMetric.WithLabelValues(autoscaler.Name, autoscaler.Namespace).Set(float64(*autoscaler.Spec.MinReplicas))
		}
		c.specMaxMetric.WithLabelValues(autoscaler.Name, autoscaler.Namespace).Set(float64(autoscaler.Spec.MaxReplicas))
		for _, metric := range autoscaler.Spec.Metrics {
			c.specMetricTarget.WithLabelValues(autoscaler.Name, autoscaler.Namespace, metric.External.Metric.Name).Set(metric.External.Target.Value.AsApproximateFloat64())
		}

		c.scaleStateDesiredMetric.WithLabelValues(autoscaler.Name, autoscaler.Namespace).Set(float64(autoscaler.Status.DesiredReplicas))
		c.scaleStateCurrentMetric.WithLabelValues(autoscaler.Name, autoscaler.Namespace).Set(float64(autoscaler.Status.CurrentReplicas))
		for _, metric := range autoscaler.Status.CurrentMetrics {
			c.metricsCurrentMetric.WithLabelValues(autoscaler.Name, autoscaler.Namespace, metric.External.Metric.Name).Set(metric.External.Current.Value.AsApproximateFloat64())
		}
	}

	c.specMinMetric.Collect(ch)
	c.specMaxMetric.Collect(ch)
	c.specMetricTarget.Collect(ch)
	c.scaleStateDesiredMetric.Collect(ch)
	c.scaleStateCurrentMetric.Collect(ch)
	c.metricsCurrentMetric.Collect(ch)
}
