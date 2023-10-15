// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func EncodeMetricHPA(hpaName string) *metav1.LabelSelector {
	return metav1.SetAsLabelSelector(labels.Set{
		"hpa": hpaName,
	})
}

func DecodeMetricHPA(selector labels.Selector) string {
	parts := strings.Split(selector.String(), "=")
	if len(parts) != 2 {
		panic(fmt.Sprintf("unexpected metric selector format: %v", selector))
	}
	if parts[0] != "hpa" {
		panic(fmt.Sprintf("unexpected metric selector label: %v", selector))
	}

	return parts[1]
}
