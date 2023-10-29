// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package metrics

import (
	"context"
	"fmt"
	"time"

	metricstypes "k9s-autoscaler/pkg/metrics/types"
	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/metrics/proto"
	scalingtypes "k9s-autoscaler/pkg/scale/types"

	protob "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

const (
	SimProviderName = "sim"
)

var (
	simSingleton *metricsSim = nil
)

type autoscalerState struct {
	totalLoadTimespan    time.Duration
	currentInstanceCount int32
}

// Simulation metrics provider adapter for testing. It also provides a scaling
// provider adapter such that when scale changes it can recaculate the load
// for metrics.
// It implements a single metric that returns the average load as a percentage.
// see pkg/providers/metrics/proto/sim.proto
type metricsSim struct {
	config                           *proto.SimMetricsConfig
	startTime                        time.Time
	autoscalerStateByNamespaceByName map[string]map[string]*autoscalerState
	initialized                      bool
}

func init() {
	providers.RegisterMetricsClient(SimProviderName, &metricsSim{})
	providers.RegisterScalingClient(SimProviderName, &metricsSim{})
}

func (s *metricsSim) MetricsClient(config *anypb.Any) (metricstypes.MetricsClient, error) {
	if s.initialized {
		return s, nil
	}

	simConfig := proto.SimMetricsConfig{}
	if err := anypb.UnmarshalTo(config, &simConfig, protob.UnmarshalOptions{}); err != nil {
		return nil, err
	}

	if len(simConfig.MetricName) == 0 {
		return nil, fmt.Errorf("metric name must be provided")
	}
	if len(simConfig.AutoscalersConfig) == 0 {
		return nil, fmt.Errorf("no autoscaler configurations specified")
	}

	s.autoscalerStateByNamespaceByName = make(map[string]map[string]*autoscalerState)
	for _, config := range simConfig.AutoscalersConfig {
		if len(config.AutoscalerName) == 0 {
			return nil, fmt.Errorf("autoscaler name cannot be empty for autoscaler %s namespace %s", config.AutoscalerName, config.AutoscalerNamespace)
		}
		if config.MaxLoadPerInstance <= 0 {
			return nil, fmt.Errorf("max load per instance must be > 0 for autoscaler %s namespace %s", config.AutoscalerName, config.AutoscalerNamespace)
		}
		if len(config.Load) == 0 {
			return nil, fmt.Errorf("no load specified to autoscaler %s namespace %s", config.AutoscalerName, config.AutoscalerNamespace)
		}
		totalTimespan := time.Duration(0)
		for _, load := range config.Load {
			if load.Timespan.AsDuration() <= time.Duration(0) {
				return nil, fmt.Errorf("timespan cannot be 0 for autoscaler %s namespace %s", config.AutoscalerName, config.AutoscalerNamespace)
			}
			totalTimespan += load.Timespan.AsDuration()
		}
		if _, ok := s.autoscalerStateByNamespaceByName[config.AutoscalerNamespace]; !ok {
			s.autoscalerStateByNamespaceByName[config.AutoscalerNamespace] = make(map[string]*autoscalerState)
		}
		s.autoscalerStateByNamespaceByName[config.AutoscalerNamespace][config.AutoscalerName] = &autoscalerState{
			totalLoadTimespan:    totalTimespan,
			currentInstanceCount: 1,
		}
	}

	s.config = &simConfig
	s.startTime = time.Now()
	s.initialized = true

	simSingleton = s

	return s, nil
}

func (s *metricsSim) ScalingClient(config *anypb.Any) (scalingtypes.ScalingClient, error) {
	if simSingleton == nil {
		return nil, fmt.Errorf("metrics sim has not been initialized or not configured")
	}

	return simSingleton, nil
}

func (s *metricsSim) GetMetric(ctx context.Context, autoscaler, namespace, metricName string) ([]int64, time.Time, error) {
	if metricName != s.config.MetricName {
		return nil, time.Time{}, fmt.Errorf("invalid metric name: %s != %s", metricName, s.config.MetricName)
	}

	for _, autoscalerConfig := range s.config.AutoscalersConfig {
		if autoscalerConfig.AutoscalerNamespace != namespace || autoscalerConfig.AutoscalerName != autoscaler {
			continue
		}
		var state *autoscalerState
		if autoscalers, ok := s.autoscalerStateByNamespaceByName[namespace]; !ok {
			return nil, time.Time{}, fmt.Errorf("no autoscaler found namespace %s", namespace)
		} else if state, ok = autoscalers[autoscaler]; !ok {
			return nil, time.Time{}, fmt.Errorf("autoscaler not found namespace %s name %s", namespace, autoscaler)
		}

		delta := time.Since(s.startTime) % state.totalLoadTimespan
		currentOffset := time.Duration(0)
		for _, load := range autoscalerConfig.Load {
			if currentOffset+load.Timespan.AsDuration() > delta {
				values := []int64{}
				// calculate load percentage
				if state.currentInstanceCount > 0 {
					value := int64(100 * (load.Load / float64(state.currentInstanceCount)) / autoscalerConfig.MaxLoadPerInstance)
					// value should be in millis
					value *= 1000
					values = append(values, value)
				}
				klog.V(10).InfoS("returning metric", "values", values)
				return values, time.Now(), nil
			}
			currentOffset += load.Timespan.AsDuration()
		}

		// shouldn't happen
		panic("expected to find load")
	}

	return nil, time.Time{}, fmt.Errorf("autoscaler %s namespace %s not found", autoscaler, namespace)
}

func (s *metricsSim) SetScaleTarget(ctx context.Context, name, namespace string, target *prototypes.ScaleSpec) error {
	var state *autoscalerState
	if autoscalers, ok := s.autoscalerStateByNamespaceByName[namespace]; !ok {
		return fmt.Errorf("no autoscaler found namespace %s", namespace)
	} else if state, ok = autoscalers[name]; !ok {
		return fmt.Errorf("autoscaler not found namespace %s name %s", namespace, name)
	}

	state.currentInstanceCount = target.Desired

	return nil
}

func (s *metricsSim) GetScale(ctx context.Context, name, namespace string) (*prototypes.Scale, error) {
	var state *autoscalerState
	if autoscalers, ok := s.autoscalerStateByNamespaceByName[namespace]; !ok {
		return nil, fmt.Errorf("no autoscaler found namespace %s", namespace)
	} else if state, ok = autoscalers[name]; !ok {
		return nil, fmt.Errorf("autoscaler not found namespace %s name %s", namespace, name)
	}

	return &prototypes.Scale{
		Spec: &prototypes.ScaleSpec{
			Desired: state.currentInstanceCount,
		},
		Status: &prototypes.ScaleStatus{
			Current: state.currentInstanceCount,
		},
	}, nil
}
