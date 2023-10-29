// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"sync"

	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	_ watch.Interface = &autoscalerWatch{}
)

type unwatcher interface {
	removeWatch(watch *autoscalerWatch)
}

// apimachinery watch adapter.
type autoscalerWatch struct {
	sync.RWMutex
	unwatcher unwatcher
	c         chan watch.Event
}

func newWatch(unwatcher unwatcher) *autoscalerWatch {
	return &autoscalerWatch{
		unwatcher: unwatcher,
		c:         make(chan watch.Event),
	}
}

// Stop stops watching. Will close the channel returned by ResultChan(). Releases
// any resources used by the watch.
func (w *autoscalerWatch) Stop() {
	w.Lock()
	defer w.Unlock()

	if w.c != nil {
		close(w.c)
		w.c = nil
		w.unwatcher.removeWatch(w)
	}
}

// ResultChan returns a chan which will receive all the events. If an error occurs
// or Stop() is called, the implementation will close this channel and
// release any resources used by the watch.
func (w *autoscalerWatch) ResultChan() <-chan watch.Event {
	return w.c
}

func (w *autoscalerWatch) add(hpa *v2.HorizontalPodAutoscaler) {
	w.c <- watch.Event{
		Type:   watch.Added,
		Object: hpa,
	}
}

func (w *autoscalerWatch) update(hpa *v2.HorizontalPodAutoscaler) {
	w.c <- watch.Event{
		Type:   watch.Modified,
		Object: hpa,
	}
}

func (w *autoscalerWatch) delete(hpa *v2.HorizontalPodAutoscaler) {
	w.c <- watch.Event{
		Type:   watch.Deleted,
		Object: hpa,
	}
}
