// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package cmd

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var RootCMD = &cobra.Command{
	Use:   "k9s-autoscaler",
	Short: "K9s autoscaler controller",
}

var (
	MetricsListenAddress string
	PProfListenAddress   string
)

func init() {
	klog.InitFlags(nil)

	RootCMD.PersistentFlags().StringVar(&MetricsListenAddress, "metrics-listen", ":8080", "prometheus metric exposer listen address")
	RootCMD.PersistentFlags().StringVar(&PProfListenAddress, "pprof-listen", ":6060", "go pprof http listen address")
}

func SetupTelemetryAndLogging() {
	setupPProf(PProfListenAddress)
	setupMetrics(MetricsListenAddress)
}

func setupPProf(pprofListenAddress string) {
	if pprofListenAddress != "" {
		go func() {
			klog.V(1).InfoS("starting pprof http handler", "listen", pprofListenAddress)
			err := http.ListenAndServe(pprofListenAddress, nil)
			klog.V(1).InfoS("pprof http handler terminated", "error", err)
		}()
	}
}

func setupMetrics(metricsListenAddress string) {
	if len(metricsListenAddress) == 0 {
		return
	}

	path := "/metrics"
	index := strings.Index(metricsListenAddress, "/")
	if index != -1 {
		path = metricsListenAddress[index:]
		metricsListenAddress = metricsListenAddress[0:index]
	}
	http.Handle(path, promhttp.Handler())
	go func() {
		klog.V(1).InfoS("starting prometheus exposer", "listen", metricsListenAddress)
		err := http.ListenAndServe(metricsListenAddress, nil)
		klog.V(1).InfoS("prometheus metrics exposer terminated", "error", err)
	}()
}
