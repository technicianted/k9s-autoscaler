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

type client struct {
	metricsclient.MetricsClient

	callbackClient types.MetricsClient
}

func NewClient(callbackClient types.MetricsClient) metricsclient.MetricsClient {
	return &client{
		callbackClient: callbackClient,
	}
}

func (c *client) GetExternalMetric(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error) {
	autoscalerName := storage.DecodeMetricHPA(selector)
	return c.callbackClient.GetMetric(context.TODO(), autoscalerName, namespace, metricName)
}
