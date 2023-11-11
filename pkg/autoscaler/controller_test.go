// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package autoscaler

import (
	"context"
	"testing"
	"time"

	"k9s-autoscaler/pkg/events"
	eventsmocks "k9s-autoscaler/pkg/events/mocks"
	"k9s-autoscaler/pkg/metrics"
	metricsmocks "k9s-autoscaler/pkg/metrics/mocks"
	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/scale"
	scalemocks "k9s-autoscaler/pkg/scale/mocks"
	"k9s-autoscaler/pkg/storage"
	storagemocks "k9s-autoscaler/pkg/storage/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/klog/v2"
)

func TestAutoscalerControllerNoScale(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	autoscalerUpdateMock := storagemocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	autoscalerUpdateMock.EXPECT().AutoscalerStatusUpdated(gomock.Any()).DoAndReturn(
		func(autoscaler *prototypes.Autoscaler) {
			klog.InfoS("status update", "status", autoscaler.Status)
		}).AnyTimes()

	storageClient, err := storage.NewClient(autoscalerUpdateMock)
	require.NoError(t, err)

	eventerMock := eventsmocks.NewMockEventCreator(mockCtrl)
	eventerMock.EXPECT().Create(gomock.Any(), t.Name(), t.Name(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, name, namespace string, event *prototypes.AutoscalerEvent) error {
			klog.InfoS("event", "event", event)
			return nil
		}).AnyTimes()
	eventNamespacer := events.NewGetter(eventerMock)

	scalerMock := scalemocks.NewMockScalingClient(mockCtrl)
	scaleGetter := scale.NewGetter(storageClient, scalerMock)
	scalerMock.EXPECT().GetScale(gomock.Any(), t.Name(), t.Name(), gomock.Any()).Return(
		&prototypes.Scale{
			Spec: &prototypes.ScaleSpec{
				Desired: 1,
			},
			Status: &prototypes.ScaleStatus{
				Current: 1,
			},
		},
		nil).AnyTimes()
	scalerMock.EXPECT().SetScaleTarget(gomock.Any(), t.Name(), t.Name(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, name, namespace string, target *prototypes.ScaleSpec) error {
			klog.InfoS("setScaleTarget", "target", target)
			return nil
		}).AnyTimes()

	metricsCallbackMock := metricsmocks.NewMockMetricsClient(mockCtrl)
	metricsCallbackMock.EXPECT().GetMetric(gomock.Any(), t.Name(), t.Name(), "testmetric").Return(
		[]int64{1.0},
		time.Now(),
		nil).AnyTimes()
	metricsGetter := metrics.NewClient(metricsCallbackMock)

	resyncPeriod := time.Second
	downscaleStabilisationWindow := 100 * time.Millisecond

	controller := NewController(
		storageClient,
		eventNamespacer,
		scaleGetter,
		metricsGetter,
		resyncPeriod,
		downscaleStabilisationWindow,
		0.1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go controller.Run(ctx, 1)

	err = storageClient.Add(&prototypes.Autoscaler{
		Name:      t.Name(),
		Namespace: t.Name(),
		Spec: &prototypes.AutoscalerSpec{
			Min: 1,
			Max: 4,
			Metrics: []*prototypes.Metric{
				{
					Name:   "testmetric",
					Target: 1.0,
				},
			},
		},
	})
	require.NoError(t, err)

	<-ctx.Done()
}
