K9s Autoscaler
---
K9s Autoscaler is a weekend project to explore the idea of running Kubernetes Horizontal Pod Autoscaler as a standalone component outside of a Kubernetes cluster, to drive autoscaling for generic external objects using external metrics.

### Objectives
* Run outside Kubernetes environment while using latest and greates code from HPA.
* Provide configuration file based autoscaler discovery and configuration.
* Build metrics backends for generic Prometheus servers and Azure Monitor.
* Build REST HTTP API to discover and configure autoscalers.

#### Project status
Early stage development. End to end unit test runs. See `pkg/autoscaler/controller_test.go`.

#### Kubernetes version
v1.27.6