// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package events

import (
	"k9s-autoscaler/pkg/events/types"

	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
)

var (
	_ v1core.EventsGetter = &getter{}
)

type getter struct {
	creator types.EventCreator
}

func NewGetter(creator types.EventCreator) v1core.EventsGetter {
	return &getter{
		creator: creator,
	}
}

func (g *getter) Events(namespace string) v1core.EventInterface {
	return NewEventer(namespace, g.creator)
}
