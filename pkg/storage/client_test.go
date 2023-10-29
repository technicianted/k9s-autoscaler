// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"context"
	"testing"
	"time"

	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/storage/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func TestClientCRUD(t *testing.T) {
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

	hpa, err := client.HorizontalPodAutoscalers(autoscaler.Namespace).Get(context.Background(), autoscaler.Name, v1.GetOptions{})
	require.NoError(t, err)
	require.NotNil(t, hpa)
	require.Equal(t, autoscaler.Name, hpa.Name)
	require.Equal(t, autoscaler.Namespace, hpa.Namespace)
	require.EqualValues(t, autoscaler.Spec.Min, *hpa.Spec.MinReplicas)

	// update
	autoscaler.Spec.Min = 10
	err = client.Update(&autoscaler)
	require.NoError(t, err)
	hpa, err = client.HorizontalPodAutoscalers(autoscaler.Namespace).Get(context.Background(), autoscaler.Name, v1.GetOptions{})
	require.NoError(t, err)
	require.EqualValues(t, autoscaler.Spec.Min, *hpa.Spec.MinReplicas)

	// delete
	err = client.Delete(autoscaler.Name, autoscaler.Namespace)
	require.NoError(t, err)
	_, err = client.HorizontalPodAutoscalers(autoscaler.Namespace).Get(context.Background(), autoscaler.Name, v1.GetOptions{})
	require.Error(t, err)
	require.True(t, errors.IsNotFound(err))
}

func TestClientK8sCacheOperationsBefore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	statusUpdateHandler := mocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	client, err := NewClient(statusUpdateHandler)
	require.NoError(t, err)

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.HorizontalPodAutoscalers("").List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.HorizontalPodAutoscalers("").Watch(context.TODO(), options)
			},
		},
		&v2.HorizontalPodAutoscaler{},
		time.Second,
		nil)

	autoscaler := prototypes.Autoscaler{
		Name:      "testas",
		Namespace: "none",
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
	// adds
	err = client.Add(&autoscaler)
	require.NoError(t, err)

	stopChan := make(chan struct{})
	go informer.Run(stopChan)
	defer close(stopChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	synced := cache.WaitForCacheSync(ctx.Done(), informer.HasSynced)
	require.True(t, synced)

	time.Sleep(100 * time.Millisecond)
	obj, exists, err := informer.GetIndexer().GetByKey(autoscaler.Namespace + "/" + autoscaler.Name)
	require.NoError(t, err)
	require.True(t, exists)
	hpa := obj.(*v2.HorizontalPodAutoscaler)
	require.EqualValues(t, autoscaler.Spec.Max, hpa.Spec.MaxReplicas)
}

func TestClientK8sCacheOperationsAfter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	statusUpdateHandler := mocks.NewMockAutoscalerStatusUpdateHandler(mockCtrl)
	client, err := NewClient(statusUpdateHandler)
	require.NoError(t, err)

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.HorizontalPodAutoscalers("").List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.HorizontalPodAutoscalers("").Watch(context.TODO(), options)
			},
		},
		&v2.HorizontalPodAutoscaler{},
		time.Second,
		nil)
	stopChan := make(chan struct{})
	go informer.Run(stopChan)
	defer close(stopChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	synced := cache.WaitForCacheSync(ctx.Done(), informer.HasSynced)
	require.True(t, synced)

	autoscaler := prototypes.Autoscaler{
		Name:      "testas",
		Namespace: "none",
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
	// adds
	err = client.Add(&autoscaler)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	obj, exists, err := informer.GetIndexer().GetByKey(autoscaler.Namespace + "/" + autoscaler.Name)
	require.NoError(t, err)
	require.True(t, exists)
	hpa := obj.(*v2.HorizontalPodAutoscaler)
	require.EqualValues(t, autoscaler.Spec.Max, hpa.Spec.MaxReplicas)

	autoscaler.Name += "2"
	err = client.Add(&autoscaler)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	obj, exists, err = informer.GetIndexer().GetByKey(autoscaler.Namespace + "/" + autoscaler.Name)
	require.NoError(t, err)
	require.True(t, exists)
	hpa2 := obj.(*v2.HorizontalPodAutoscaler)
	require.EqualValues(t, autoscaler.Spec.Max, hpa2.Spec.MaxReplicas)

	var hpas []*v2.HorizontalPodAutoscaler
	cache.ListAll(informer.GetIndexer(), labels.Everything(), func(m interface{}) {
		hpas = append(hpas, m.(*v2.HorizontalPodAutoscaler))
	})
	require.Len(t, hpas, 2)

	// updates
	autoscaler.Spec.Max = 10
	err = client.Update(&autoscaler)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	obj, exists, err = informer.GetIndexer().GetByKey(autoscaler.Namespace + "/" + autoscaler.Name)
	require.NoError(t, err)
	require.True(t, exists)
	hpa2 = obj.(*v2.HorizontalPodAutoscaler)
	require.EqualValues(t, autoscaler.Spec.Max, hpa2.Spec.MaxReplicas)

	// delete
	err = client.Delete(autoscaler.Name, autoscaler.Namespace)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	_, exists, err = informer.GetIndexer().GetByKey(autoscaler.Namespace + "/" + autoscaler.Name)
	require.NoError(t, err)
	require.False(t, exists)
}
