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

func (am *azureMonitor) GetMetric(ctx context.Context, metricName, namespace string, selector map[string]string) ([]int64, time.Time, error) {
	options := azquery.MetricsClientQueryResourceOptions{
		MetricNames: &metricName,
		Interval:    to.Ptr("PT1M"),
		Timespan:    to.Ptr(azquery.NewTimeInterval(time.Now().Add(-10*time.Minute), time.Now())),
	}

	resourceURI, ok := selector[proto.AzureMonitorConfig_ResourceURI.String()]
	if !ok {
		return nil, time.Time{}, fmt.Errorf("%s selector is required", proto.AzureMonitorConfig_ResourceURI.String())
	}
	if v, ok := selector[proto.AzureMonitorConfig_MetricNamespace.String()]; !ok {
		return nil, time.Time{}, fmt.Errorf("%s selector is required", proto.AzureMonitorConfig_MetricNamespace.String())
	} else {
		options.MetricNamespace = &v
	}
	if v, ok := selector[proto.AzureMonitorConfig_Filter.String()]; ok {
		options.Filter = &v
	}

	response, err := am.metricsClient.QueryResource(
		ctx,
		resourceURI,
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

	return []int64{int64(*latest.Average)}, *latest.TimeStamp, nil
}

func (f *azureMonitorFactory) MetricsClient(config *anypb.Any) (metricstypes.MetricsClient, error) {
	return newAzureMonitor()
}
