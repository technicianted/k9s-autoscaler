K9s Autoscaler
---
K9s Autoscaler is a weekend project to explore the idea of running Kubernetes Horizontal Pod Autoscaler as a standalone component outside of a Kubernetes cluster, to drive autoscaling for generic external objects using external metrics.

### Objectives
* Run outside Kubernetes environment while using latest and greates code from HPA.
* Provide configuration file based autoscaler discovery and configuration.
* Build metrics backends for generic Prometheus servers and Azure Monitor.
* Build REST HTTP API to discover and configure autoscalers.
* Fully extensible in and out of tree for storage, metrics and scaling adapters.

#### Project status
Early stage development.

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
        # list of metrics and their targets. metricsClient must be able to provide
        # values for these.
        - name: testmetric
          target: 70
# metricsClient provides an adapter for reading metrics values.
metricsClient:
  # sim is a metricsClient and scalingClient that can be used to simulate scaling
  # and metrics reading based on predefined time intervals.
  name: sim
  config:
    "@type": type.googleapis.com/k9sautoscaler.providers.metrics.proto.SimMetricsConfig
    metricName: testmetric
    autoscalersConfig:
    - autoscalerName: testauto1
      autoscalerNamespace: testnamespace
      maxLoadPerInstance: 10.0
      load:
      - timespan: 10s
        load: 100
      - timespan: 10s
        load: 200 
      - timespan: 10s
        load: 50
# scalingClient provides an adapter for getting and setting scale.
scalingClient:
  # sim is a scalingClient that works with sim metrics client.
  name: sim
# eventsClient provides an adapter for the autoscaler events.
eventsClient:
  # klog logs status updates to klogger.
  name: klog
resyncPeriod: 5s
```

#### Kubernetes version
v1.27.6