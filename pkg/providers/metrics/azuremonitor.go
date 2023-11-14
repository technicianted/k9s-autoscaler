// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"fmt"
	"time"

	metricstypes "k9s-autoscaler/pkg/metrics/types"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/metrics/proto"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	protob "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type azureMonitor struct {
	metricsClient *azquery.MetricsClient
}

type azureMonitorFactory struct {
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

	options := azquery.MetricsClientQueryResourceOptions{
		MetricNames: &metricName,
		Interval:    to.Ptr("PT1M"),
		Timespan:    to.Ptr(azquery.NewTimeInterval(time.Now().Add(-10*time.Minute), time.Now())),
	}

	if len(metricConfig.ResourceURI) == 0 {
		return nil, time.Time{}, fmt.Errorf("resourceURI is required")
	}
	if len(metricConfig.MetricNamespace) == 0 {
		return nil, time.Time{}, fmt.Errorf("namespace is required")
	}
	options.MetricNamespace = &metricConfig.MetricNamespace
	if metricConfig.Filter != nil {
		options.Filter = metricConfig.Filter
	}

	response, err := am.metricsClient.QueryResource(
		ctx,
		metricConfig.ResourceURI,
		&options)
	if err != nil {
		return nil, time.Time{}, err
	}

	if len(response.Value) == 0 {
		return nil, time.Now(), fmt.Errorf("no metric values available")
	}
	ts := response.Value[0].TimeSeries
	if len(ts) != 1 {
		return nil, time.Now(), fmt.Errorf("expecting 1 timeseries, got %d", len(ts))
	}
	values := ts[0].Data
	if len(values) == 0 {
		return nil, time.Now(), fmt.Errorf("no timeseries data")
	}
	latest := values[len(values)-1]

	var value *float64
	switch metricConfig.Aggregation {
	case proto.AzureMonitorMetricConfig_None:
		return nil, time.Time{}, fmt.Errorf("aggregation type is required")
	case proto.AzureMonitorMetricConfig_Average:
		value = latest.Average
	case proto.AzureMonitorMetricConfig_Count:
		value = latest.Count
	case proto.AzureMonitorMetricConfig_Maximum:
		value = latest.Maximum
	case proto.AzureMonitorMetricConfig_Minimum:
		value = latest.Minimum
	case proto.AzureMonitorMetricConfig_Total:
		value = latest.Total
	case proto.AzureMonitorMetricConfig_RatePerMinute:
		if latest.Total == nil {
			return nil, time.Time{}, fmt.Errorf("cannot calculate rate: metric does not support Total")
		}
		value = to.Ptr(*latest.Total / float64(60))
	default:
		return nil, time.Time{}, fmt.Errorf("unknown aggregation type: %v", metricConfig.Aggregation)
	}
	if value == nil {
		return nil, time.Time{}, fmt.Errorf("specified aggregation type %v is not supported by metric", metricConfig.Aggregation.String())
	}

	return []int64{int64(*value)}, *latest.TimeStamp, nil
}

func (f *azureMonitorFactory) MetricsClient(config *anypb.Any) (metricstypes.MetricsClient, error) {
	return newAzureMonitor()
}
