// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package main

import (
	"fmt"
	"os"

	"k9s-autoscaler/cmd/k9s-autoscaler/cmd"
)

func main() {
	if err := cmd.RootCMD.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
