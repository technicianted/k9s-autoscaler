// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"context"
	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/storage/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func TestNamespacedClientWatches(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	statusUpdateHandler := mocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	client, err := NewClient(statusUpdateHandler)
	require.NoError(t, err)

	autoscaler := prototypes.Autoscaler{
		Name:      "testas",
		Namespace: "testasns",
		Spec: &prototypes.AutoscalerSpec{
			Min: 1,
			Max: 2,
			Metrics: []*prototypes.Metric{
				{
					Name:   "testmetric",
					Target: 1,
				},
			},
		},
	}

	w, err := client.HorizontalPodAutoscalers(autoscaler.Namespace).Watch(context.Background(), v1.ListOptions{})
	require.NoError(t, err)
	defer w.Stop()

	watcherDoneChan := make(chan bool)
	go func() {
		e := <-w.ResultChan()
		assert.Equal(t, watch.Added, e.Type)
		hpa := e.Object.(*v2.HorizontalPodAutoscaler)
		require.Equal(t, autoscaler.Name, hpa.Name)
		watcherDoneChan <- true

		e = <-w.ResultChan()
		assert.Equal(t, watch.Modified, e.Type)
		hpa = e.Object.(*v2.HorizontalPodAutoscaler)
		require.Equal(t, autoscaler.Name, hpa.Name)
		require.EqualValues(t, autoscaler.Spec.Min, *hpa.Spec.MinReplicas)
		watcherDoneChan <- true

		e = <-w.ResultChan()
		assert.Equal(t, watch.Deleted, e.Type)
		hpa = e.Object.(*v2.HorizontalPodAutoscaler)
		require.Equal(t, autoscaler.Name, hpa.Name)
		watcherDoneChan <- true
	}()

	err = client.Add(&autoscaler)
	require.NoError(t, err)
	select {
	case <-time.After(time.Second):
		require.Fail(t, "timed out waiting for watcher add")
	case <-watcherDoneChan:
	}

	autoscaler.Spec.Min = 10
	err = client.Update(&autoscaler)
	require.NoError(t, err)
	select {
	case <-time.After(time.Second):
		require.Fail(t, "timed out waiting for watcher update")
	case <-watcherDoneChan:
	}

	autoscaler.Spec.Min = 10
	err = client.Delete(autoscaler.Name, autoscaler.Namespace)
	require.NoError(t, err)
	select {
	case <-time.After(time.Second):
		require.Fail(t, "timed out waiting for watcher delete")
	case <-watcherDoneChan:
	}
}

func TestNamespacedClientWatchesAllNamespaces(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	statusUpdateHandler := mocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	client, err := NewClient(statusUpdateHandler)
	require.NoError(t, err)

	autoscaler := prototypes.Autoscaler{
		Name:      "testas",
		Namespace: "testasns",
		Spec: &prototypes.AutoscalerSpec{
			Min: 1,
			Max: 2,
			Metrics: []*prototypes.Metric{
				{
					Name:   "testmetric",
					Target: 1,
				},
			},
		},
	}

	w, err := client.HorizontalPodAutoscalers(v1.NamespaceAll).Watch(context.Background(), v1.ListOptions{})
	require.NoError(t, err)
	defer w.Stop()

	watcherDoneChan := make(chan bool)
	go func() {
		e := <-w.ResultChan()
		assert.Equal(t, watch.Added, e.Type)
		hpa := e.Object.(*v2.HorizontalPodAutoscaler)
		require.Equal(t, autoscaler.Name, hpa.Name)
		watcherDoneChan <- true

		e = <-w.ResultChan()
		assert.Equal(t, watch.Modified, e.Type)
		hpa = e.Object.(*v2.HorizontalPodAutoscaler)
		require.Equal(t, autoscaler.Name, hpa.Name)
		require.EqualValues(t, autoscaler.Spec.Min, *hpa.Spec.MinReplicas)
		watcherDoneChan <- true

		e = <-w.ResultChan()
		assert.Equal(t, watch.Deleted, e.Type)
		hpa = e.Object.(*v2.HorizontalPodAutoscaler)
		require.Equal(t, autoscaler.Name, hpa.Name)
		watcherDoneChan <- true
	}()

	err = client.Add(&autoscaler)
	require.NoError(t, err)
	select {
	case <-time.After(time.Second):
		require.Fail(t, "timed out waiting for watcher add")
	case <-watcherDoneChan:
	}

	autoscaler.Spec.Min = 10
	err = client.Update(&autoscaler)
	require.NoError(t, err)
	select {
	case <-time.After(time.Second):
		require.Fail(t, "timed out waiting for watcher update")
	case <-watcherDoneChan:
	}

	autoscaler.Spec.Min = 10
	err = client.Delete(autoscaler.Name, autoscaler.Namespace)
	require.NoError(t, err)
	select {
	case <-time.After(time.Second):
		require.Fail(t, "timed out waiting for watcher delete")
	case <-watcherDoneChan:
	}
}

func TestNamespacedClientUpdateStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	statusUpdateHandler := mocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	client, err := NewClient(statusUpdateHandler)
	require.NoError(t, err)

	autoscaler := prototypes.Autoscaler{
		Name:      "testas",
		Namespace: "testasns",
		Spec: &prototypes.AutoscalerSpec{
			Min: 1,
			Max: 2,
			Metrics: []*prototypes.Metric{
				{
					Name:   "testmetric",
					Target: 1,
				},
			},
		},
	}

	err = client.Add(&autoscaler)
	require.NoError(t, err)

	statusUpdateHandler.EXPECT().AutoscalerStatusUpdated(gomock.Any()).DoAndReturn(
		func(as *prototypes.Autoscaler) {
			assert.Equal(t, int32(10), *as.Status.CurrentScale)
		})

	updatedHPA, err := client.HorizontalPodAutoscalers(autoscaler.Namespace).UpdateStatus(
		context.Background(),
		&v2.HorizontalPodAutoscaler{
			ObjectMeta: v1.ObjectMeta{
				Name:      autoscaler.Name,
				Namespace: autoscaler.Namespace,
			},
			Status: v2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 10,
			},
		},
		v1.UpdateOptions{})
	require.NoError(t, err)
	require.EqualValues(t, 10, updatedHPA.Status.CurrentReplicas)
}

func TestNamespacedClientErrors(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	statusUpdateHandler := mocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	client, err := NewClient(statusUpdateHandler)
	require.NoError(t, err)

	// get
	_, err = client.HorizontalPodAutoscalers("none").Get(context.Background(), "none", v1.GetOptions{})
	require.Error(t, err)
	require.True(t, errors.IsNotFound(err))

	// update
	// invalid namespace in object
	_, err = client.HorizontalPodAutoscalers("none").UpdateStatus(
		context.Background(),
		&v2.HorizontalPodAutoscaler{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "none",
			},
		},
		v1.UpdateOptions{})
	require.Error(t, err)
	require.True(t, errors.IsNotFound(err))

	_, err = client.HorizontalPodAutoscalers("none").UpdateStatus(
		context.Background(),
		&v2.HorizontalPodAutoscaler{
			ObjectMeta: v1.ObjectMeta{},
		},
		v1.UpdateOptions{})
	require.Error(t, err)
	require.True(t, errors.IsConflict(err))
}
