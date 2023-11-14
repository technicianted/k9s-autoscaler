// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"fmt"
	"time"

	"k9s-autoscaler/pkg/metrics/types"
	"k9s-autoscaler/pkg/storage"
	storagetypes "k9s-autoscaler/pkg/storage/types"

	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/labels"
	metricsclient "k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
)

// Adapter partial implemenation of k8s MetricsClient. Since k9s autoscaler
// uses k8s autoscaler ExternalMetrics, only these functions are implemented.
type client struct {
	metricsclient.MetricsClient

	autoscalerGetter storagetypes.AutoscalerGetter
	callbackClient   types.MetricsClient
}

// Create a new adapter k8s MetricsClient that calls callbackClient to get metric
// values.
func NewClient(autoscalerGetter storagetypes.AutoscalerGetter, callbackClient types.MetricsClient) metricsclient.MetricsClient {
	return &client{
		autoscalerGetter: autoscalerGetter,
		callbackClient:   callbackClient,
	}
}

func (c *client) GetExternalMetric(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error) {
	autoscalerName := storage.DecodeMetricHPA(selector)
	as, err := c.autoscalerGetter.Get(autoscalerName, namespace)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to get autoscaler: %v", err)
	}

	var config *anypb.Any
	for _, metric := range as.Spec.Metrics {
		if metric.Name == metricName {
			config = metric.Config
			break
		}
	}
	if config == nil {
		return nil, time.Time{}, fmt.Errorf("could not find configs for metric %s", metricName)
	}

	startTime := time.Now()
	values, ts, err := c.callbackClient.GetMetric(
		context.TODO(),
		metricName,
		autoscalerName,
		namespace,
		config)
	if err != nil {
		metricLatencyMetric.WithLabelValues(namespace, metricName, "true").Observe(float64(time.Since(startTime)))
		return nil, time.Time{}, err
	}
	metricLatencyMetric.WithLabelValues(namespace, metricName, "").Observe(float64(time.Since(startTime)))

	// autoscaler expect millis representation
	for i, value := range values {
		values[i] = value * 1000
	}

	return values, ts, nil
}
