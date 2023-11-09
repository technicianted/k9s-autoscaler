// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k9s-autoscaler/pkg/metrics/types"

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
	mapSelector, err := convertSelectorToLabelsMap(selector)
	if err != nil {
		return nil, time.Time{}, err
	}

	startTime := time.Now()
	values, ts, err := c.callbackClient.GetMetric(
		context.TODO(),
		metricName,
		namespace, mapSelector)
	if err != nil {
		metricLatencyMetric.WithLabelValues(namespace, metricName, "true").Observe(float64(time.Since(startTime)))
		return nil, time.Time{}, err
	}
	metricLatencyMetric.WithLabelValues(namespace, metricName, "").Observe(float64(time.Since(startTime)))

	return values, ts, nil
}

func convertSelectorToLabelsMap(selector labels.Selector) (map[string]string, error) {
	selectorMap := map[string]string{}
	exprs := strings.Split(selector.String(), ",")
	for _, expr := range exprs {
		parts := strings.Split(expr, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid selector expression: %s", expr)
		}
		selectorMap[parts[0]] = parts[1]
	}

	return selectorMap, nil
}
