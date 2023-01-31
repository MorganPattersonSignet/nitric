// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/nitric/core/pkg/worker/pool (interfaces: WorkerPool)

// Package worker is a generated GoMock package.
package worker

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	worker "github.com/nitrictech/nitric/core/pkg/worker"
	pool "github.com/nitrictech/nitric/core/pkg/worker/pool"
)

// MockWorkerPool is a mock of WorkerPool interface.
type MockWorkerPool struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerPoolMockRecorder
}

// MockWorkerPoolMockRecorder is the mock recorder for MockWorkerPool.
type MockWorkerPoolMockRecorder struct {
	mock *MockWorkerPool
}

// NewMockWorkerPool creates a new mock instance.
func NewMockWorkerPool(ctrl *gomock.Controller) *MockWorkerPool {
	mock := &MockWorkerPool{ctrl: ctrl}
	mock.recorder = &MockWorkerPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkerPool) EXPECT() *MockWorkerPoolMockRecorder {
	return m.recorder
}

// AddWorker mocks base method.
func (m *MockWorkerPool) AddWorker(arg0 worker.Worker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWorker", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddWorker indicates an expected call of AddWorker.
func (mr *MockWorkerPoolMockRecorder) AddWorker(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWorker", reflect.TypeOf((*MockWorkerPool)(nil).AddWorker), arg0)
}

// GetWorker mocks base method.
func (m *MockWorkerPool) GetWorker(arg0 *pool.GetWorkerOptions) (worker.Worker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorker", arg0)
	ret0, _ := ret[0].(worker.Worker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorker indicates an expected call of GetWorker.
func (mr *MockWorkerPoolMockRecorder) GetWorker(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorker", reflect.TypeOf((*MockWorkerPool)(nil).GetWorker), arg0)
}

// GetWorkerCount mocks base method.
func (m *MockWorkerPool) GetWorkerCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkerCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetWorkerCount indicates an expected call of GetWorkerCount.
func (mr *MockWorkerPoolMockRecorder) GetWorkerCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkerCount", reflect.TypeOf((*MockWorkerPool)(nil).GetWorkerCount))
}

// GetWorkers mocks base method.
func (m *MockWorkerPool) GetWorkers(arg0 *pool.GetWorkerOptions) []worker.Worker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkers", arg0)
	ret0, _ := ret[0].([]worker.Worker)
	return ret0
}

// GetWorkers indicates an expected call of GetWorkers.
func (mr *MockWorkerPoolMockRecorder) GetWorkers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkers", reflect.TypeOf((*MockWorkerPool)(nil).GetWorkers), arg0)
}

// Monitor mocks base method.
func (m *MockWorkerPool) Monitor() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Monitor")
	ret0, _ := ret[0].(error)
	return ret0
}

// Monitor indicates an expected call of Monitor.
func (mr *MockWorkerPoolMockRecorder) Monitor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Monitor", reflect.TypeOf((*MockWorkerPool)(nil).Monitor))
}

// RemoveWorker mocks base method.
func (m *MockWorkerPool) RemoveWorker(arg0 worker.Worker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveWorker", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveWorker indicates an expected call of RemoveWorker.
func (mr *MockWorkerPoolMockRecorder) RemoveWorker(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveWorker", reflect.TypeOf((*MockWorkerPool)(nil).RemoveWorker), arg0)
}

// WaitForMinimumWorkers mocks base method.
func (m *MockWorkerPool) WaitForMinimumWorkers(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForMinimumWorkers", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForMinimumWorkers indicates an expected call of WaitForMinimumWorkers.
func (mr *MockWorkerPoolMockRecorder) WaitForMinimumWorkers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForMinimumWorkers", reflect.TypeOf((*MockWorkerPool)(nil).WaitForMinimumWorkers), arg0)
}