// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/nitric/core/pkg/workers/websockets (interfaces: WebsocketRequestHandler)

// Package mock_websockets is a generated GoMock package.
package mock_websockets

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

// MockWebsocketRequestHandler is a mock of WebsocketRequestHandler interface.
type MockWebsocketRequestHandler struct {
	ctrl     *gomock.Controller
	recorder *MockWebsocketRequestHandlerMockRecorder
}

// MockWebsocketRequestHandlerMockRecorder is the mock recorder for MockWebsocketRequestHandler.
type MockWebsocketRequestHandlerMockRecorder struct {
	mock *MockWebsocketRequestHandler
}

// NewMockWebsocketRequestHandler creates a new mock instance.
func NewMockWebsocketRequestHandler(ctrl *gomock.Controller) *MockWebsocketRequestHandler {
	mock := &MockWebsocketRequestHandler{ctrl: ctrl}
	mock.recorder = &MockWebsocketRequestHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebsocketRequestHandler) EXPECT() *MockWebsocketRequestHandlerMockRecorder {
	return m.recorder
}

// HandleRequest mocks base method.
func (m *MockWebsocketRequestHandler) HandleRequest(arg0 *websocketspb.ServerMessage) (*websocketspb.ClientMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleRequest", arg0)
	ret0, _ := ret[0].(*websocketspb.ClientMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleRequest indicates an expected call of HandleRequest.
func (mr *MockWebsocketRequestHandlerMockRecorder) HandleRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleRequest", reflect.TypeOf((*MockWebsocketRequestHandler)(nil).HandleRequest), arg0)
}

// WorkerCount mocks base method.
func (m *MockWebsocketRequestHandler) WorkerCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkerCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// WorkerCount indicates an expected call of WorkerCount.
func (mr *MockWebsocketRequestHandlerMockRecorder) WorkerCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkerCount", reflect.TypeOf((*MockWebsocketRequestHandler)(nil).WorkerCount))
}