// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package cmd

import (
	controllercmd "k9s-autoscaler/pkg/cmd"
	"k9s-autoscaler/pkg/version"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var ControllerCMD = &cobra.Command{
	Use:   "controller",
	Short: "K9s autoscaler controller",

	Run: runController,
}

var (
	opts = controllercmd.NewOptions()
)

func init() {
	ControllerCMD.Flags().StringVar(&opts.YAMLConfigPath, "config", opts.YAMLConfigPath, "path to yaml configuration file")
	ControllerCMD.MarkFlagFilename("config")
	ControllerCMD.MarkFlagRequired("config")
	ControllerCMD.Flags().IntVar(&opts.Workers, "workers", opts.Workers, "number of controller workers")
	RootCMD.AddCommand(ControllerCMD)
}

func runController(command *cobra.Command, args []string) {
	SetupTelemetryAndLogging()
	klog.InfoS("starting k9s autoscaler controller", "version", version.Build)

	c, err := controllercmd.NewControllerCMD(opts)
	if err != nil {
		klog.Exitf("failed to start controller: %v", err)
	}
	go c.Start()

	klog.InfoS("startup sequence completed")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	klog.InfoS("shutting down")
	c.Stop()
}
