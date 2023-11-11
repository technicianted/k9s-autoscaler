// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"testing"
	"time"

	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers/metrics/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestMetricsSimSimple(t *testing.T) {
	config := proto.SimMetricsConfig{
		MetricName: t.Name(),
		AutoscalersConfig: []*proto.AutoscalerConfig{
			{
				AutoscalerName:      t.Name(),
				AutoscalerNamespace: "testnamespace",
				MaxLoadPerInstance:  50,
				Load: []*proto.MetricLoad{
					{
						Timespan: durationpb.New(100 * time.Millisecond),
						Load:     100,
					},
					{
						Timespan: durationpb.New(200 * time.Millisecond),
						Load:     200,
					},
				},
			},
		},
	}
	configAny, err := anypb.New(&config)
	require.NoError(t, err)
	client, err := (&metricsSim{}).MetricsClient(configAny)
	require.NoError(t, err)
	sim := client.(*metricsSim)

	selector := map[string]string{
		proto.SimMetricsConfig_AUTOSCALER_NAME.String(): t.Name(),
	}

	err = sim.SetScaleTarget(
		context.Background(),
		t.Name(),
		"testnamespace",
		&prototypes.AutoscalerTarget{},
		&prototypes.ScaleSpec{Desired: 1})
	require.NoError(t, err)
	values, _, err := client.GetMetric(context.Background(), t.Name(), "testnamespace", selector)
	require.NoError(t, err)
	require.Len(t, values, 1)
	// 200%
	require.EqualValues(t, 200000, values[0])
	time.Sleep(100 * time.Millisecond)
	values, _, err = client.GetMetric(context.Background(), t.Name(), "testnamespace", selector)
	require.NoError(t, err)
	require.Len(t, values, 1)
	// 400%
	require.EqualValues(t, 400000, values[0])

	// back from the start
	time.Sleep(200 * time.Millisecond)
	values, _, err = client.GetMetric(context.Background(), t.Name(), "testnamespace", selector)
	require.NoError(t, err)
	require.Len(t, values, 1)
	// 200%
	require.EqualValues(t, 200000, values[0])
}
