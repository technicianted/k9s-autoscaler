// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"context"
	"time"

	v2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	autoscalinginformers "k8s.io/client-go/informers/autoscaling/v2"
	v2listers "k8s.io/client-go/listers/autoscaling/v2"
	"k8s.io/client-go/tools/cache"
)

var (
	_ autoscalinginformers.HorizontalPodAutoscalerInformer = &HPAInformer{}
)

type HPAInformer struct {
	hpaInformer cache.SharedIndexInformer
}

func NewInformer(storageClient *Client) *HPAInformer {
	informer := &HPAInformer{
		hpaInformer: cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
					return storageClient.HorizontalPodAutoscalers(v1.NamespaceAll).List(context.TODO(), options)
				},
				WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
					return storageClient.HorizontalPodAutoscalers(v1.NamespaceAll).Watch(context.TODO(), options)
				},
			},
			&v2.HorizontalPodAutoscaler{},
			time.Second,
			nil),
	}

	return informer
}

func (i *HPAInformer) Run(stopCh <-chan struct{}) {
	go i.hpaInformer.Run(stopCh)
}

func (i *HPAInformer) Informer() cache.SharedIndexInformer {
	return i.hpaInformer
}

func (i *HPAInformer) Lister() v2listers.HorizontalPodAutoscalerLister {
	return v2listers.NewHorizontalPodAutoscalerLister(i.hpaInformer.GetIndexer())
}
