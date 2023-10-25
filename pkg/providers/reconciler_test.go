// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package providers

import (
	"testing"

	prototypes "k9s-autoscaler/pkg/proto"
	storagemocks "k9s-autoscaler/pkg/storage/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestReconcilerSimpleAddUpdateDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := storagemocks.NewMockAutoscalerCRUDder(mockCtrl)

	list := []*prototypes.Autoscaler{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
	}
	client.EXPECT().List().Return([]*prototypes.Autoscaler{}, nil)
	client.EXPECT().Add(list[0]).Return(nil)
	client.EXPECT().Add(list[1]).Return(nil)

	r := NewReconciler(client)
	err := r.Reconcile(list)
	assert.NoError(t, err)

	// update
	updatedList := []*prototypes.Autoscaler{
		proto.Clone(list[0]).(*prototypes.Autoscaler),
		proto.Clone(list[1]).(*prototypes.Autoscaler),
	}
	updatedList[0].Spec = &prototypes.AutoscalerSpec{Max: 1}
	client.EXPECT().List().Return(list, nil)
	client.EXPECT().Update(updatedList[0])
	err = r.Reconcile(updatedList)
	assert.NoError(t, err)

	// delete
	updatedList = updatedList[1:]
	client.EXPECT().List().Return(list, nil)
	client.EXPECT().Delete(list[0].Name, list[0].Namespace)
	err = r.Reconcile(updatedList)
	assert.NoError(t, err)
}
