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

	"google.golang.org/protobuf/proto"
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

	scalingClients                      = make(map[string]scalingtypes.ScalingClient)
	scalingClientFactoryByTargetConfigs = make(map[string]string)
)

// Registers a storage client provider adapter with configMEssage and factory.
// configMessage type will be used for identification.
func RegisterStorageClient(configMessage proto.Message, factory StorageClientFactory) {
	name := typeNameForMessage(configMessage)
	if _, ok := storageClientFactories[name]; ok {
		panic(fmt.Sprintf("storage provider %s already registered", name))
	}

	storageClientFactories[name] = factory
}

// Instantiates a new storage client with config. Name in config is used to lookup
// previously registerd factory.
func StorageClient(config *configproto.ProviderConfig) (*storage.Client, error) {
	name := config.Config.TypeUrl
	fmt.Printf("type: %+v\n", config.Config)
	if f, ok := storageClientFactories[name]; !ok {
		return nil, fmt.Errorf("no storage client registered for %s", name)
	} else {
		return f.StorageClient(config.Config)
	}
}

// Registers a storage client provider adapter with name and factory.
func RegisterMetricsClient(configMessage proto.Message, factory MetricsClientFactory) {
	name := typeNameForMessage(configMessage)
	if _, ok := metricClientFactories[name]; ok {
		panic(fmt.Sprintf("metrics client %s already registered", name))
	}

	metricClientFactories[name] = factory
}

// Get a storage client with config. Name in config is used to lookup
// previously registerd factory.
func MetricsClient(config *configproto.ProviderConfig) (metricstypes.MetricsClient, error) {
	name := config.Config.TypeUrl
	if f, ok := metricClientFactories[name]; !ok {
		return nil, fmt.Errorf("no metrics client registered for %s", name)
	} else {
		return f.MetricsClient(config.Config)
	}
}

// Registers a storage client provider adapter with config and targetConfig and factory.
func RegisterScalingClient(configMessage proto.Message, targetConfigMessage proto.Message, factory ScalingClientFactory) {
	name := typeNameForMessage(configMessage)
	if _, ok := scalingClientFactories[name]; ok {
		panic(fmt.Sprintf("scaling client %s already registered", name))
	}
	targetName := typeNameForMessage(targetConfigMessage)
	if _, ok := scalingClientFactoryByTargetConfigs[targetName]; ok {
		panic(fmt.Sprintf("scaling client target config %s already registered", targetName))
	}
	scalingClientFactories[name] = factory
	scalingClientFactoryByTargetConfigs[targetName] = name
}

// Gets a scaling client with config. Name in config is used to lookup
// previously registerd factory.
func ScalingClient(config *configproto.ProviderConfig) (scalingtypes.ScalingClient, error) {
	name := config.Config.TypeUrl
	if f, ok := scalingClientFactories[name]; !ok {
		return nil, fmt.Errorf("no scaling client registered for %s", name)
	} else {
		return f.ScalingClient(config.Config)
	}
}

// Gets a scaling client from a scaling target configuration.
func ScalingClientByTargetConfig(config *anypb.Any) (scalingtypes.ScalingClient, error) {
	name := config.TypeUrl
	clientName, ok := scalingClientFactoryByTargetConfigs[name]
	if !ok {
		return nil, fmt.Errorf("no scaling client registered for %s", name)
	}
	client, ok := scalingClients[clientName]
	if !ok {
		return nil, fmt.Errorf("client %s for target %s not configured", clientName, name)
	}

	return client, nil
}

// Registers an events client provider adapter with name and factory.
func RegisterEventsClient(configMessage proto.Message, factory EventsClientFactory) {
	name := typeNameForMessage(configMessage)
	if _, ok := eventsClientFactories[name]; ok {
		panic(fmt.Sprintf("events client %s already registered", name))
	}

	eventsClientFactories[name] = factory
}

// Gets an events client with config. Name in config is used to lookup
// previously registerd factory.
func EventsClient(config *configproto.ProviderConfig) (eventstypes.EventCreator, error) {
	name := config.Config.TypeUrl
	if f, ok := eventsClientFactories[name]; !ok {
		return nil, fmt.Errorf("no events client registered for %s", name)
	} else {
		return f.EventsClient(config.Config)
	}
}

func typeNameForMessage(configMessage proto.Message) string {
	return "type.googleapis.com/" + string(configMessage.ProtoReflect().Descriptor().FullName())
}
