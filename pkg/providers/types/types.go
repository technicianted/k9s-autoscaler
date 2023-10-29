// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import (
	prototypes "k9s-autoscaler/pkg/proto"
)

//go:generate mockgen -package mocks -destination ../mocks/reconciler.go -source $GOFILE

// A simple implementation that can reconcile two lists of autoscalers.
type Reconciler interface {
	// Reconcile attempts to reconcile autoscalers into storage. Returns a wrapped
	// joined errors with failures encountered.
	Reconcile(autoscalers []*prototypes.Autoscaler) error
}
