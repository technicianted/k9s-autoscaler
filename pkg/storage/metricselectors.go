// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package storage

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Encodes hpaName into k8s LabelSelector. Since k9s API exposes metrics as
// part of autoscalers, call adapter calls require identification of its name
// and namespace. This function uses LabelSelector to encode such information.
func EncodeMetricHPA(hpaName string) *metav1.LabelSelector {
	return metav1.SetAsLabelSelector(labels.Set{
		"hpa": hpaName,
	})
}

// Decodes hpaName into k8s LabelSelector. Since k9s API exposes metrics as
// part of autoscalers, call adapter calls require identification of its name
// and namespace. This function uses LabelSelector to encode such information.
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
