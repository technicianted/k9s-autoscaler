// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	proto "k9s-autoscaler/pkg/proto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockScalingClient is a mock of ScalingClient interface.
type MockScalingClient struct {
	ctrl     *gomock.Controller
	recorder *MockScalingClientMockRecorder
}

// MockScalingClientMockRecorder is the mock recorder for MockScalingClient.
type MockScalingClientMockRecorder struct {
	mock *MockScalingClient
}

// NewMockScalingClient creates a new mock instance.
func NewMockScalingClient(ctrl *gomock.Controller) *MockScalingClient {
	mock := &MockScalingClient{ctrl: ctrl}
	mock.recorder = &MockScalingClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScalingClient) EXPECT() *MockScalingClientMockRecorder {
	return m.recorder
}

// GetScale mocks base method.
func (m *MockScalingClient) GetScale(ctx context.Context, name, namespace string, scaleTarget *proto.AutoscalerTarget) (*proto.Scale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScale", ctx, name, namespace, scaleTarget)
	ret0, _ := ret[0].(*proto.Scale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScale indicates an expected call of GetScale.
func (mr *MockScalingClientMockRecorder) GetScale(ctx, name, namespace, scaleTarget interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScale", reflect.TypeOf((*MockScalingClient)(nil).GetScale), ctx, name, namespace, scaleTarget)
}

// SetScaleTarget mocks base method.
func (m *MockScalingClient) SetScaleTarget(ctx context.Context, name, namespace string, scaleTarget *proto.AutoscalerTarget, target *proto.ScaleSpec) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetScaleTarget", ctx, name, namespace, scaleTarget, target)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetScaleTarget indicates an expected call of SetScaleTarget.
func (mr *MockScalingClientMockRecorder) SetScaleTarget(ctx, name, namespace, scaleTarget, target interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetScaleTarget", reflect.TypeOf((*MockScalingClient)(nil).SetScaleTarget), ctx, name, namespace, scaleTarget, target)
}
