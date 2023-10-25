// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package types

import "context"

type Controller interface {
	Run(ctx context.Context, workers int)
}
