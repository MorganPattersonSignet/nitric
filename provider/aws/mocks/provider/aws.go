// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/nitric/provider/aws/runtime/core (interfaces: AwsProvider)

// Package mock_core is a generated GoMock package.
package mock_core

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/nitrictech/nitric/core/pkg/providers/common"
)

// MockAwsProvider is a mock of AwsProvider interface.
type MockAwsProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAwsProviderMockRecorder
}

// MockAwsProviderMockRecorder is the mock recorder for MockAwsProvider.
type MockAwsProviderMockRecorder struct {
	mock *MockAwsProvider
}

// NewMockAwsProvider creates a new mock instance.
func NewMockAwsProvider(ctrl *gomock.Controller) *MockAwsProvider {
	mock := &MockAwsProvider{ctrl: ctrl}
	mock.recorder = &MockAwsProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAwsProvider) EXPECT() *MockAwsProviderMockRecorder {
	return m.recorder
}

// Details mocks base method.
func (m *MockAwsProvider) Details(arg0 context.Context, arg1, arg2 string) (*common.DetailsResponse[interface{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Details", arg0, arg1, arg2)
	ret0, _ := ret[0].(*common.DetailsResponse[interface{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Details indicates an expected call of Details.
func (mr *MockAwsProviderMockRecorder) Details(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Details", reflect.TypeOf((*MockAwsProvider)(nil).Details), arg0, arg1, arg2)
}

// GetResources mocks base method.
func (m *MockAwsProvider) GetResources(arg0 context.Context, arg1 string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResources", arg0, arg1)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResources indicates an expected call of GetResources.
func (mr *MockAwsProviderMockRecorder) GetResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResources", reflect.TypeOf((*MockAwsProvider)(nil).GetResources), arg0, arg1)
}