// Code generated by MockGen. DO NOT EDIT.
// Source: services/queue.go
//
// Generated by this command:
//
//	mockgen -source services/queue.go -destination services/queue_mock.go
//

// Package mock_services is a generated GoMock package.
package services

import (
	types "container-manager/types"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockQueue is a mock of Queue interface.
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue.
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance.
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Enqueue mocks base method.
func (m *MockQueue) Enqueue(jobID string, container types.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", jobID, container)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockQueueMockRecorder) Enqueue(jobID, container any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockQueue)(nil).Enqueue), jobID, container)
}

// GetStatus mocks base method.
func (m *MockQueue) GetStatus(jobID string) (types.JobStatus, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", jobID)
	ret0, _ := ret[0].(types.JobStatus)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockQueueMockRecorder) GetStatus(jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockQueue)(nil).GetStatus), jobID)
}

// Run mocks base method.
func (m *MockQueue) Run(workerCount int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", workerCount)
}

// Run indicates an expected call of Run.
func (mr *MockQueueMockRecorder) Run(workerCount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockQueue)(nil).Run), workerCount)
}

// Stop mocks base method.
func (m *MockQueue) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockQueueMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockQueue)(nil).Stop))
}