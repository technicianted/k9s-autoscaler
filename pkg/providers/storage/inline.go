// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"fmt"
	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/storage/proto"
	"k9s-autoscaler/pkg/storage"

	protob "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	ProviderName = "inline"
)

type inlineStorage struct{}

func init() {
	providers.RegisterStorageClient(ProviderName, &inlineStorage{})
}

// Inline storage provider defines a provider that reads autoscaler configurations
// embedded in its own convigurations.
// Once loaded, configurations are immutable.
// see: examples/intree/sim.yaml for an example.
func (f *inlineStorage) StorageClient(config *anypb.Any) (*storage.Client, error) {
	inlineConfig := proto.InlineStorageConfig{}
	if err := anypb.UnmarshalTo(config, &inlineConfig, protob.UnmarshalOptions{}); err != nil {
		return nil, err
	}
	if len(inlineConfig.Autoscalers) == 0 {
		return nil, fmt.Errorf("no autoscalers specified")
	}

	client, err := storage.NewClient(f)
	if err != nil {
		return nil, err
	}
	for _, autoscalerConfig := range inlineConfig.Autoscalers {
		if err := client.Add(autoscalerConfig); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (f *inlineStorage) AutoscalerStatusUpdated(autoscaler *prototypes.Autoscaler) {
}
