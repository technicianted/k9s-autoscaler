// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestControllerCMD(t *testing.T) {
	optionsYAML := `
storageClient:
  name: inline
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.storage.proto.InlineStorageConfig
    autoscalers:
    - name: testauto1
      namespace: testnamespace
      spec:
        min: 1
        max: 30
        metrics:
        - name: testmetric
          target: 70
metricsClient:
  name: sim
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.metrics.proto.SimMetricsConfig
    metricName: "testmetric"
    autoscalersConfig:
    - autoscalerName: testauto1
      autoscalerNamespace: testnamespace
      maxLoadPerInstance: 10.0
      load:
      - timespan: 5s
        load: 100
      - timespan: 5s
        load: 200 
      - timespan: 5s
        load: 50 
scalingClient:
  name: sim
eventsClient:
  name: klog
resyncPeriod: 1s
`
	configFile, err := os.CreateTemp("/tmp", t.Name())
	require.NoError(t, err)
	_, err = configFile.WriteString(optionsYAML)
	require.NoError(t, err)
	configFile.Close()
	defer os.Remove(configFile.Name())

	opts := NewOptions()
	opts.YAMLConfigPath = configFile.Name()
	c, err := NewControllerCMD(opts)
	require.NoError(t, err)
	go c.Start()
	time.Sleep(15 * time.Second)
	c.Stop()
}
