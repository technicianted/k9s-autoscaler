// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package autoscaler

import (
	"context"
	"time"

	"k9s-autoscaler/pkg/scale"
	"k9s-autoscaler/pkg/storage"

	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	autoscalinginformers "k8s.io/client-go/informers/autoscaling/v2"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	scaleclient "k8s.io/client-go/scale"
	"k8s.io/kubernetes/pkg/controller/podautoscaler"
	metricsclient "k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
)

type controller struct {
	apimeta.RESTMapper
	autoscalinginformers.HorizontalPodAutoscalerInformer

	storageClient *storage.Client
	hpaInformer   *storage.HPAInformer
	k8sController *podautoscaler.HorizontalController
}

func NewController(
	storageClient *storage.Client,
	evtNamespacer v1core.EventsGetter,
	scaleNamespacer scaleclient.ScalesGetter,
	metricsClient metricsclient.MetricsClient,
	resyncPeriod time.Duration,
	downscaleStabilisationWindow time.Duration,
	tolerance float64) *controller {

	c := &controller{
		storageClient: storageClient,
		hpaInformer:   storage.NewInformer(storageClient),
	}

	// do nothing pod informer
	podInformer := &podInformer{}

	mapper := c
	// unused options
	cpuInitializationPeriod := 5 * time.Minute
	delayOfInitialReadinessStatus := 30 * time.Second
	containerResourceMetricsEnabled := false

	c.k8sController = podautoscaler.NewHorizontalController(
		evtNamespacer,
		scaleNamespacer,
		storageClient,
		mapper,
		metricsClient,
		c.hpaInformer,
		podInformer,
		resyncPeriod,
		downscaleStabilisationWindow,
		tolerance,
		cpuInitializationPeriod,
		delayOfInitialReadinessStatus,
		containerResourceMetricsEnabled)

	return c
}

func (c *controller) Run(ctx context.Context, workers int) {
	go c.hpaInformer.Run(ctx.Done())

	c.k8sController.Run(ctx, workers)
}

func (c *controller) RESTMappings(gk schema.GroupKind, versions ...string) ([]*apimeta.RESTMapping, error) {
	return []*apimeta.RESTMapping{
		{
			Resource: schema.GroupVersionResource{
				Group:   scale.ScalingResourceGroup,
				Version: scale.ScalingResourceVersion,
			},
		},
	}, nil
}
