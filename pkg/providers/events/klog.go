// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package events

import (
	"context"

	eventstypes "k9s-autoscaler/pkg/events/types"
	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/events/proto"

	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

const (
	KLogProviderName = "klog"
)

// A simple event provider adapter that emits the event to klog.
type klogClient struct{}

func init() {
	providers.RegisterEventsClient(&proto.KLog{}, &klogClient{})
}

func (e *klogClient) EventsClient(config *anypb.Any) (eventstypes.EventCreator, error) {
	return e, nil
}

func (e *klogClient) Create(ctx context.Context, name, namespace string, event *prototypes.AutoscalerEvent) error {
	klog.V(0).InfoS("status update", "name", name, "namespace", namespace, event)

	return nil
}
