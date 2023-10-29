// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"time"

	"k9s-autoscaler/pkg/metrics/types"
	"k9s-autoscaler/pkg/storage"

	"k8s.io/apimachinery/pkg/labels"
	metricsclient "k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
)

// Adapter partial implemenation of k8s MetricsClient. Since k9s autoscaler
// uses k8s autoscaler ExternalMetrics, only these functions are implemented.
type client struct {
	metricsclient.MetricsClient

	callbackClient types.MetricsClient
}

// Create a new adapter k8s MetricsClient that calls callbackClient to get metric
// values.
func NewClient(callbackClient types.MetricsClient) metricsclient.MetricsClient {
	return &client{
		callbackClient: callbackClient,
	}
}

func (c *client) GetExternalMetric(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error) {
	autoscalerName := storage.DecodeMetricHPA(selector)
	startTime := time.Now()
	values, ts, err := c.callbackClient.GetMetric(context.TODO(), autoscalerName, namespace, metricName)
	if err != nil {
		metricLatencyMetric.WithLabelValues(autoscalerName, namespace, metricName, "true").Observe(float64(time.Since(startTime)))
		return nil, time.Time{}, err
	}
	metricLatencyMetric.WithLabelValues(autoscalerName, namespace, metricName, "").Observe(float64(time.Since(startTime)))

	return values, ts, nil
}
