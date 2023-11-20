// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	metricstypes "k9s-autoscaler/pkg/metrics/types"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/metrics/proto"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	protob "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

type azureMonitor struct {
	metricsClient *azquery.MetricsClient
}

type azureMonitorFactory struct {
}

type azureMonitorMetric struct {
	name        string
	aggregation proto.AzureMonitorMetricConfig_Aggregation
}
type azureMonitorMetricGroup struct {
	metrics     map[string]azureMonitorMetric
	namespace   string
	resourceURI string
	filter      string
}
type azureMonitorTimeSeriesResult struct {
	dimensions map[string]string
	value      int64
}

func init() {
	providers.RegisterMetricsClient(&proto.AzureMonitorConfig{}, &azureMonitorFactory{})
}

func newAzureMonitor() (*azureMonitor, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get default azure credentials: %v", err)
	}
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azure monitor metrics client: %v", err)
	}

	return &azureMonitor{
		metricsClient: client,
	}, nil
}

func (am *azureMonitor) GetMetric(ctx context.Context, metricName, autoscalerName, namespace string, config *anypb.Any) ([]int64, time.Time, error) {
	metricConfig := proto.AzureMonitorMetricConfig{}
	if err := anypb.UnmarshalTo(config, &metricConfig, protob.UnmarshalOptions{}); err != nil {
		return nil, time.Time{}, err
	}

	metrics := azureMonitorMetricGroup{
		metrics: map[string]azureMonitorMetric{
			metricName: {
				name:        metricName,
				aggregation: metricConfig.Aggregation,
			},
		},
		namespace:   metricConfig.MetricNamespace,
		resourceURI: metricConfig.ResourceURI,
	}
	if metricConfig.Filter != nil {
		metrics.filter = *metricConfig.Filter
	}

	results, timestamp, err := getMetricValues(ctx, am.metricsClient, metrics)
	if err != nil {
		return nil, time.Time{}, err
	}
	if len(results) != 1 {
		return nil, time.Time{}, fmt.Errorf("expecting 1 metric, got %d", len(results))
	}
	timeseries, ok := results[metricName]
	if !ok {
		return nil, time.Time{}, fmt.Errorf("did not get expected metric in results")
	}
	if len(timeseries) != 1 {
		return nil, time.Time{}, fmt.Errorf("expecting 1 timeseries, got %d", len(timeseries))
	}

	klog.InfoS("azure monitor metrics", "values", results, "timestamp", timestamp)

	return []int64{timeseries[0].value}, timestamp, nil
}

func (f *azureMonitorFactory) MetricsClient(config *anypb.Any) (metricstypes.MetricsClient, error) {
	return newAzureMonitor()
}

func getMetricValues(ctx context.Context, metricsClient *azquery.MetricsClient, metrics azureMonitorMetricGroup) (map[string][]azureMonitorTimeSeriesResult, time.Time, error) {
	if len(metrics.resourceURI) == 0 {
		return nil, time.Time{}, fmt.Errorf("resourceURI is required")
	}
	if len(metrics.namespace) == 0 {
		return nil, time.Time{}, fmt.Errorf("namespace is required")
	}

	metricNames := []string{}
	for name := range metrics.metrics {
		metricNames = append(metricNames, name)
	}

	options := azquery.MetricsClientQueryResourceOptions{
		MetricNames: to.Ptr(strings.Join(metricNames, ",")),
		Interval:    to.Ptr("PT1M"),
		Timespan:    to.Ptr(azquery.NewTimeInterval(time.Now().Add(-10*time.Minute), time.Now())),
	}

	options.MetricNamespace = &metrics.namespace
	if len(metrics.filter) > 0 {
		options.Filter = &metrics.filter
	}

	response, err := metricsClient.QueryResource(
		ctx,
		metrics.resourceURI,
		&options)
	if err != nil {
		return nil, time.Time{}, err
	}
	if len(response.Value) != len(metrics.metrics) {
		return nil, time.Now(), fmt.Errorf("mismatched number of returned metrics: %d != %d", len(metrics.metrics), len(response.Value))
	}

	metricsValues := map[string][]azureMonitorTimeSeriesResult{}
	// TODO: this assumes that latest is same across all returned metrics
	latestTimestamp := time.Time{}
	for _, value := range response.Value {
		incomingName := *value.Name.Value
		for _, ts := range value.TimeSeries {
			timeSeriesResult := azureMonitorTimeSeriesResult{dimensions: make(map[string]string)}
			for _, dim := range ts.MetadataValues {
				timeSeriesResult.dimensions[*dim.Name.Value] = *dim.Value
			}

			values := ts.Data
			if len(values) == 0 {
				return nil, time.Now(), fmt.Errorf("no timeseries data for %s", incomingName)
			}
			latest := values[len(values)-1]

			var metricValue *float64
			switch metrics.metrics[incomingName].aggregation {
			case proto.AzureMonitorMetricConfig_None:
				return nil, time.Time{}, fmt.Errorf("aggregation type is required")
			case proto.AzureMonitorMetricConfig_Average:
				metricValue = latest.Average
			case proto.AzureMonitorMetricConfig_Count:
				metricValue = latest.Count
			case proto.AzureMonitorMetricConfig_Maximum:
				metricValue = latest.Maximum
			case proto.AzureMonitorMetricConfig_Minimum:
				metricValue = latest.Minimum
			case proto.AzureMonitorMetricConfig_Total:
				metricValue = latest.Total
			case proto.AzureMonitorMetricConfig_RatePerMinute:
				if latest.Total == nil {
					return nil, time.Time{}, fmt.Errorf("cannot calculate rate: metric does not support Total")
				}
				metricValue = latest.Total
			default:
				return nil, time.Time{}, fmt.Errorf("unknown aggregation type: %v", metrics.metrics[incomingName].aggregation)
			}
			if metricValue == nil {
				return nil, time.Time{}, fmt.Errorf("specified aggregation type %v is not supported by metric", metrics.metrics[incomingName].aggregation.String())
			}
			latestTimestamp = *latest.TimeStamp
			timeSeriesResult.value = int64(*metricValue)
			metricsValues[incomingName] = append(metricsValues[incomingName], timeSeriesResult)
		}
	}

	return metricsValues, latestTimestamp, nil
}
