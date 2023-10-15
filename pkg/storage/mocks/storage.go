// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	proto "k9s-autoscaler/pkg/proto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAutoscalerCRUDder is a mock of AutoscalerCRUDder interface.
type MockAutoscalerCRUDder struct {
	ctrl     *gomock.Controller
	recorder *MockAutoscalerCRUDderMockRecorder
}

// MockAutoscalerCRUDderMockRecorder is the mock recorder for MockAutoscalerCRUDder.
type MockAutoscalerCRUDderMockRecorder struct {
	mock *MockAutoscalerCRUDder
}

// NewMockAutoscalerCRUDder creates a new mock instance.
func NewMockAutoscalerCRUDder(ctrl *gomock.Controller) *MockAutoscalerCRUDder {
	mock := &MockAutoscalerCRUDder{ctrl: ctrl}
	mock.recorder = &MockAutoscalerCRUDderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutoscalerCRUDder) EXPECT() *MockAutoscalerCRUDderMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockAutoscalerCRUDder) Add(autoscaler *proto.Autoscaler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", autoscaler)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockAutoscalerCRUDderMockRecorder) Add(autoscaler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAutoscalerCRUDder)(nil).Add), autoscaler)
}

// Delete mocks base method.
func (m *MockAutoscalerCRUDder) Delete(name, namespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", name, namespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAutoscalerCRUDderMockRecorder) Delete(name, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAutoscalerCRUDder)(nil).Delete), name, namespace)
}

// Update mocks base method.
func (m *MockAutoscalerCRUDder) Update(autoscaler *proto.Autoscaler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", autoscaler)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAutoscalerCRUDderMockRecorder) Update(autoscaler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAutoscalerCRUDder)(nil).Update), autoscaler)
}

// MockAutoscalerStatusUpdateHandler is a mock of AutoscalerStatusUpdateHandler interface.
type MockAutoscalerStatusUpdateHandler struct {
	ctrl     *gomock.Controller
	recorder *MockAutoscalerStatusUpdateHandlerMockRecorder
}

// MockAutoscalerStatusUpdateHandlerMockRecorder is the mock recorder for MockAutoscalerStatusUpdateHandler.
type MockAutoscalerStatusUpdateHandlerMockRecorder struct {
	mock *MockAutoscalerStatusUpdateHandler
}

// NewMockAutoscalerStatusUpdateHandler creates a new mock instance.
func NewMockAutoscalerStatusUpdateHandler(ctrl *gomock.Controller) *MockAutoscalerStatusUpdateHandler {
	mock := &MockAutoscalerStatusUpdateHandler{ctrl: ctrl}
	mock.recorder = &MockAutoscalerStatusUpdateHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutoscalerStatusUpdateHandler) EXPECT() *MockAutoscalerStatusUpdateHandlerMockRecorder {
	return m.recorder
}

// AutoscalerStatusUpdated mocks base method.
func (m *MockAutoscalerStatusUpdateHandler) AutoscalerStatusUpdated(autoscaler *proto.Autoscaler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AutoscalerStatusUpdated", autoscaler)
}

// AutoscalerStatusUpdated indicates an expected call of AutoscalerStatusUpdated.
func (mr *MockAutoscalerStatusUpdateHandlerMockRecorder) AutoscalerStatusUpdated(autoscaler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoscalerStatusUpdated", reflect.TypeOf((*MockAutoscalerStatusUpdateHandler)(nil).AutoscalerStatusUpdated), autoscaler)
}
