// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"k9s-autoscaler/pkg/autoscaler"
	autoscalertypes "k9s-autoscaler/pkg/autoscaler/types"
	configproto "k9s-autoscaler/pkg/cmd/proto"
	"k9s-autoscaler/pkg/events"
	eventstypes "k9s-autoscaler/pkg/events/types"
	"k9s-autoscaler/pkg/metrics"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/scale"

	_ "k9s-autoscaler/pkg/providers/events"
	_ "k9s-autoscaler/pkg/providers/metrics"
	_ "k9s-autoscaler/pkg/providers/storage"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"sigs.k8s.io/yaml"
)

// Define a core autoscaler controller command abstraction that can be used
// to construct concrete instances. For example using a CLI.
type ControllerCMD struct {
	opts       Options
	controller autoscalertypes.Controller
	cancel     context.CancelFunc
}

// Creates a new instanec with opts. Returned controller command must be started
// by calling Start().
func NewControllerCMD(opts Options) (*ControllerCMD, error) {
	if len(opts.YAMLConfigPath) == 0 {
		return nil, fmt.Errorf("config path must be specified")
	}

	bytes, err := os.ReadFile(opts.YAMLConfigPath)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		return nil, err
	}

	configs := configproto.ControllerConfig{
		ResyncPeriod:                 durationpb.New(15 * time.Second),
		Tolerance:                    0.1,
		DownscaleStabilizationWindow: durationpb.New(5 * time.Minute),
	}
	if err = protojson.Unmarshal(jsonBytes, &configs); err != nil {
		return nil, err
	}

	controller, err := NewControllerFromConfigs(&configs)
	if err != nil {
		return nil, err
	}

	return &ControllerCMD{
		opts:       opts,
		controller: controller,
	}, nil
}

func (c *ControllerCMD) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.controller.Run(ctx, c.opts.Workers)

	return nil
}

func (c *ControllerCMD) Stop() error {
	c.cancel()

	return nil
}

// Utility function that creates a new autoscaler controller from given configuration
// proto message. It handles all the required validation and initialization of
// provider adapters.
func NewControllerFromConfigs(configs *configproto.ControllerConfig) (autoscalertypes.Controller, error) {
	if configs.StorageClient == nil {
		return nil, fmt.Errorf("no storage client specified")
	}
	if configs.MetricsClient == nil {
		return nil, fmt.Errorf("no metrics client specified")
	}
	if configs.ScalingClient == nil {
		return nil, fmt.Errorf("no scaling client specified")
	}

	storageClient, err := providers.StorageClient(configs.StorageClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}
	metricsClient, err := providers.MetricsClient(configs.MetricsClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %v", err)
	}
	scalingClient, err := providers.ScalingClient(configs.ScalingClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create scaling client: %v", err)
	}
	var eventsCreator eventstypes.EventCreator
	if configs.EventsClient != nil {
		eventsCreator, err = providers.EventsClient(configs.EventsClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create events client: %v", err)
		}
	}

	controller := autoscaler.NewController(
		storageClient,
		events.NewGetter(eventsCreator),
		scale.NewGetter(storageClient, scalingClient),
		metrics.NewClient(metricsClient),
		configs.ResyncPeriod.AsDuration(),
		configs.DownscaleStabilizationWindow.AsDuration(),
		configs.Tolerance)

	return controller, nil
}
