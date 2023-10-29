// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package providers

import (
	"fmt"

	eventstypes "k9s-autoscaler/pkg/events/types"
	metricstypes "k9s-autoscaler/pkg/metrics/types"
	configproto "k9s-autoscaler/pkg/providers/proto"
	scalingtypes "k9s-autoscaler/pkg/scale/types"
	"k9s-autoscaler/pkg/storage"

	"google.golang.org/protobuf/types/known/anypb"
)

// Defines a registration of a factory of storage clients.
type StorageClientFactory interface {
	// Create a new storage client from config. Implementations must validate
	// config according to their own proto scheams.
	StorageClient(config *anypb.Any) (*storage.Client, error)
}

// Defines a registration of a factory of metrics clients.
type MetricsClientFactory interface {
	// Create a new metrics from config. Implementations must validate
	// config according to their own proto scheams.
	MetricsClient(config *anypb.Any) (metricstypes.MetricsClient, error)
}

// Defines a registration of a factory of scaling clients.
type ScalingClientFactory interface {
	// Create a new scaling client from config. Implementations must validate
	// config according to their own proto scheams.
	ScalingClient(config *anypb.Any) (scalingtypes.ScalingClient, error)
}

// Defines a registration of a factory of events clients.
type EventsClientFactory interface {
	// Create a new events client from config. Implementations must validate
	// config according to their own proto scheams.
	EventsClient(config *anypb.Any) (eventstypes.EventCreator, error)
}

var (
	storageClientFactories = make(map[string]StorageClientFactory)
	metricClientFactories  = make(map[string]MetricsClientFactory)
	scalingClientFactories = make(map[string]ScalingClientFactory)
	eventsClientFactories  = make(map[string]EventsClientFactory)
)

// Registers a storage client provider adapter with name and factory.
func RegisterStorageClient(name string, factory StorageClientFactory) {
	if _, ok := storageClientFactories[name]; ok {
		panic(fmt.Sprintf("storage client %s already registered", name))
	}

	storageClientFactories[name] = factory
}

// Instantiates a new storage client with config. Name in config is used to lookup
// previously registerd factory.
func StorageClient(config *configproto.ProviderConfig) (*storage.Client, error) {
	if f, ok := storageClientFactories[config.Name]; !ok {
		return nil, fmt.Errorf("no storage client registered for %s", config.Name)
	} else {
		return f.StorageClient(config.Config)
	}
}

// Registers a storage client provider adapter with name and factory.
func RegisterMetricsClient(name string, factory MetricsClientFactory) {
	if _, ok := metricClientFactories[name]; ok {
		panic(fmt.Sprintf("metrics client %s already registered", name))
	}

	metricClientFactories[name] = factory
}

// Get a storage client with config. Name in config is used to lookup
// previously registerd factory.
func MetricsClient(config *configproto.ProviderConfig) (metricstypes.MetricsClient, error) {
	if f, ok := metricClientFactories[config.Name]; !ok {
		return nil, fmt.Errorf("no metrics client registered for %s", config.Name)
	} else {
		return f.MetricsClient(config.Config)
	}
}

// Registers a storage client provider adapter with name and factory.
func RegisterScalingClient(name string, factory ScalingClientFactory) {
	if _, ok := scalingClientFactories[name]; ok {
		panic(fmt.Sprintf("scaling client %s already registered", name))
	}

	scalingClientFactories[name] = factory
}

// Gets a scaling client with config. Name in config is used to lookup
// previously registerd factory.
func ScalingClient(config *configproto.ProviderConfig) (scalingtypes.ScalingClient, error) {
	if f, ok := scalingClientFactories[config.Name]; !ok {
		return nil, fmt.Errorf("no scaling client registered for %s", config.Name)
	} else {
		return f.ScalingClient(config.Config)
	}
}

// Registers an events client provider adapter with name and factory.
func RegisterEventsClient(name string, factory EventsClientFactory) {
	if _, ok := eventsClientFactories[name]; ok {
		panic(fmt.Sprintf("events client %s already registered", name))
	}

	eventsClientFactories[name] = factory
}

// Gets an events client with config. Name in config is used to lookup
// previously registerd factory.
func EventsClient(config *configproto.ProviderConfig) (eventstypes.EventCreator, error) {
	if f, ok := eventsClientFactories[config.Name]; !ok {
		return nil, fmt.Errorf("no events client registered for %s", config.Name)
	} else {
		return f.EventsClient(config.Config)
	}
}
