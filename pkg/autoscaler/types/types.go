// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import "context"

// A simple interface to abstract an autoscaler controller.
type Controller interface {
	// Start the controller with workers. Stop of ctx is cancelled.
	Run(ctx context.Context, workers int)
}
