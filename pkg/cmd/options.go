// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package cmd

// Struct of required options to start a ControllerCMD. Callers can provide
// the information based on their respective implmentation. For example, for
// CLI they can be provided by command line options.
type Options struct {
	// Path to configurations in yaml format.
	// see: pkg/cmd/proto/config.proto
	YAMLConfigPath string
	// Number of controller workers, defaults to 1.
	Workers int
}

// Creates new Options initialized with defaults.
func NewOptions() Options {
	return Options{
		Workers: 1,
	}
}
