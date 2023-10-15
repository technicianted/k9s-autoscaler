// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package autoscaler

import (
	"k9s-autoscaler/pkg/scale"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type podInformer struct{}

type noopSharedIndexInformer struct {
	cache.SharedIndexInformer
}

type podFakeLister struct {
	v1listers.PodLister
}

type podFakeNamespaceLister struct {
	v1listers.PodNamespaceLister

	fakeLister *podFakeLister
}

func (pi *podInformer) Informer() cache.SharedIndexInformer {
	return &noopSharedIndexInformer{}
}

func (pi *podInformer) Lister() v1listers.PodLister {
	return &podFakeLister{}
}

func (noop *noopSharedIndexInformer) HasSynced() bool {
	return true
}

func (fake *podFakeLister) List(selector labels.Selector) (ret []*v1.Pod, err error) {
	replicas := scale.DecodePodLabels(selector)
	for i := 0; i < int(replicas); i++ {
		ret = append(ret, &v1.Pod{
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
				Conditions: []v1.PodCondition{
					{
						Type:   v1.ContainersReady,
						Status: v1.ConditionTrue,
					},
				},
			},
		})
	}

	return
}

func (fake *podFakeLister) Pods(namespace string) v1listers.PodNamespaceLister {
	return &podFakeNamespaceLister{
		fakeLister: fake,
	}
}

func (fakeNamespace *podFakeNamespaceLister) List(selector labels.Selector) (ret []*v1.Pod, err error) {
	return fakeNamespace.fakeLister.List(selector)
}
