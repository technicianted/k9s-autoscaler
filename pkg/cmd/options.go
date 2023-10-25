// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package cmd

type Options struct {
	YAMLConfigPath string
	Workers        int
}

func NewOptions() Options {
	return Options{
		Workers: 1,
	}
}
