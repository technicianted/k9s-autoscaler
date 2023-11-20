package metrics

import (
	"context"
	"fmt"
	"time"

	metricstypes "k9s-autoscaler/pkg/metrics/types"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/metrics/proto"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	protob "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

// A metrics provider with a collection of ready to use Azure OpenAI metrics.
// See pkg/providers/metrics/proto/aoai.proto for more details.
type aoai struct {
	metricsClient *azquery.MetricsClient
}

type aoaiFactory struct {
}

func init() {
	providers.RegisterMetricsClient(&proto.AzureOAIConfig{}, &aoaiFactory{})
}

func newAzureOAI(config *anypb.Any) (*aoai, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get default azure credentials: %v", err)
	}
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azure monitor metrics client: %v", err)
	}

	return &aoai{
		metricsClient: client,
	}, nil
}

func (a *aoai) GetMetric(ctx context.Context, metricName, autoscalerName, namespace string, config *anypb.Any) ([]int64, time.Time, error) {
	metricConfig := proto.AzureOAIMetricConfig{}
	if err := anypb.UnmarshalTo(config, &metricConfig, protob.UnmarshalOptions{}); err != nil {
		return nil, time.Time{}, err
	}

	metricEnumValue, ok := proto.AzureOAIMetricConfig_Metric_value[metricName]
	if !ok {
		return nil, time.Time{}, fmt.Errorf("unknown metric: %s", metricName)
	}
	// only 429 is supported
	if metricEnumValue != int32(proto.AzureOAIMetricConfig_Percent429Rate) {
		return nil, time.Time{}, fmt.Errorf("unsuppoerted metric: %s", metricName)
	}
	metrics := azureMonitorMetricGroup{
		metrics: map[string]azureMonitorMetric{
			"AzureOpenAIRequests": {
				name:        "AzureOpenAIRequests",
				aggregation: proto.AzureMonitorMetricConfig_Total,
			},
		},
		namespace:   "microsoft.cognitiveservices/accounts",
		resourceURI: metricConfig.ResourceURI,
		filter:      fmt.Sprintf("ModelDeploymentName eq '%s' and StatusCode eq '*'", metricConfig.DeploymentName),
	}
	results, timestamp, err := getMetricValues(ctx, a.metricsClient, metrics)
	if err != nil {
		return nil, time.Time{}, err
	}
	if len(results) != 1 {
		return nil, time.Time{}, fmt.Errorf("expecting 1 metric, got %d", len(results))
	}
	timeseries, ok := results["AzureOpenAIRequests"]
	if !ok {
		return nil, time.Time{}, fmt.Errorf("did not get expected metric in results")
	}
	if len(timeseries) == 1 {
		return []int64{}, time.Time{}, nil
	}

	total := int64(0)
	count429 := int64(0)
	for _, ts := range timeseries {
		total += ts.value
		if ts.dimensions["statuscode"] == "429" {
			count429 += ts.value
		}
	}
	percentage := int64(float64(100*count429) / float64(total))

	klog.InfoS("aoai metric", "value", percentage)

	return []int64{percentage}, timestamp, nil
}

func (f *aoaiFactory) MetricsClient(config *anypb.Any) (metricstypes.MetricsClient, error) {
	return newAzureOAI(config)
}
