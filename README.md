K9s Autoscaler
---
K9s Autoscaler is a weekend project to explore the idea of running Kubernetes Horizontal Pod Autoscaler as a standalone component outside of a Kubernetes cluster, to drive autoscaling for generic external objects using external metrics.

### Objectives
* Run outside Kubernetes environment while using latest and greates code from HPA.
* Provide configuration file based autoscaler discovery and configuration.
* Build metrics backends for generic Prometheus servers and Azure Monitor.
* Build REST HTTP API to discover and configure autoscalers.
* Fully extensible in and out of tree for storage, metrics and scaling adapters.

### Project status
Early stage development.

#### Available storage clients
* **[Inline](pkg/providers/storage/proto/inline.proto)**: Load autoscaler configurations from a yaml file.

#### Available metrics clients
* **[Sim](pkg/providers/metrics/proto/sim.proto)**: Simulation of dummy metrics for testing.
* **[Azure Monitor](pkg/providers/metrics/proto/azuremonitor.proto)**: Read metric values of a resource from Azure Monitor metrics API.

#### Available scalers
* **[Sim](pkg/providers/metrics/proto/sim.proto)**: Simulation of dummy scaling that works with Sim metrics clients to provide proportional scale metrics.
* **[Azure Cognitive Services](pkg/providers/scaling/proto/azuredeployment.proto)**: Scales an Azure Cognitive Services resource targetting a specific deployment.

### Current usage

Command line tool is available that runs in-tree clients:
```
$ make
$ bin/k9s-autoscaler controller --config examples/intree/sim.yaml
```

Sample configuration:
```yaml
# storageClient provides an adapter for autoscaler discovery and configuration.
storageClient:
  # inline is a built-in storage client that uses in-line yaml configuration
  # for autoscalers.
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.storage.proto.InlineStorageConfig
    autoscalers:
    - name: testauto1
      namespace: testnamespace
      spec:
        min: 1
        max: 30
        target:
          # set provider target to scaling sim
          config:
            "@type": type.googleapis.com/k9sautoscaler.providers.metrics.proto.SimScalingTargetConfig
        metrics:
        # list of metrics and their targets. metricsClient must be able to provide
        # values for these.
        - name: testmetric
          target: 70
          # set metrics provider for this metric to sim
          config:
            "@type": type.googleapis.com/k9sautoscaler.providers.metrics.proto.SimMetricConfig
# metricsClient provides an adapter for reading metrics values.
metricsClient:
  # sim is a metricsClient and scalingClient that can be used to simulate scaling
  # and metrics reading based on predefined time intervals.
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.metrics.proto.SimConfig
    metricName: testmetric
    autoscalersConfig:
    - autoscalerName: testauto1
      autoscalerNamespace: testnamespace
      maxLoadPerInstance: 10.0
      load:
      - timespan: 20s
        load: 100
      - timespan: 20s
        load: 200 
      - timespan: 20s
        load: 50
# scalingClient provides an adapter for getting and setting scale.
scalingClient:
  # sim is a scalingClient that works with sim metrics client.
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.metrics.proto.SimConfig
# eventsClient provides an adapter for the autoscaler events.
eventsClient:
  # klog logs status updates to klogger.
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.events.proto.KLog
resyncPeriod: 5s
```

After running the above for a few minutes, prometheus would display something like this:

<p align="center">
  <img width="512" src="images/prom-sample-current-scale.png"/>
</p>
<p align="center">
  <img width="512" src="images/prom-sample-metrics-current.png"/>
</p>

#### Kubernetes version
v1.27.6