// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package events

import (
	"context"

	"k9s-autoscaler/pkg/events/types"
	prototypes "k9s-autoscaler/pkg/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
)

// A partial implementation of k8s EventInterface that passes Create() calls
// to the adapter callback.
type eventer struct {
	v1core.EventInterface

	creator   types.EventCreator
	namespace string
}

// Create a new k8s EventInterface for namespace that uses creator as a callback
// to create the events.
func NewEventer(namespace string, creator types.EventCreator) v1core.EventInterface {
	return &eventer{
		namespace: namespace,
		creator:   creator,
	}
}

func (e *eventer) Create(ctx context.Context, event *v1.Event, opts metav1.CreateOptions) (*v1.Event, error) {
	if e.creator != nil {
		return event, e.creator.Create(
			ctx,
			event.InvolvedObject.Name,
			e.namespace,
			eventToAutoscalerEvent(event))
	}

	return event, nil
}

// CreateWithEventNamespace is the same as a Create, except that it sends the request to the event.Namespace.
func (e *eventer) CreateWithEventNamespace(event *v1.Event) (*v1.Event, error) {
	if e.creator != nil {
		return event, e.creator.Create(
			context.TODO(),
			event.InvolvedObject.Name,
			event.Namespace,
			eventToAutoscalerEvent(event))
	}

	return event, nil
}

// PatchWithEventNamespace is the same as a Patch, except that it sends the request to the event.Namespace.
func (e *eventer) PatchWithEventNamespace(event *v1.Event, data []byte) (*v1.Event, error) {
	return e.CreateWithEventNamespace(event)
}

func eventToAutoscalerEvent(event *v1.Event) *prototypes.AutoscalerEvent {
	return &prototypes.AutoscalerEvent{
		Reason:         event.Reason,
		Message:        event.Message,
		FirstTimestamp: timestamppb.New(event.CreationTimestamp.Time),
		LastTimestamp:  timestamppb.New(event.LastTimestamp.Time),
		Count:          event.Count,
		Type:           event.Type,
		EventTime:      timestamppb.New(event.EventTime.Time),
		Action:         event.Action,
	}
}
