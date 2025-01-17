// Code generated by MockGen. DO NOT EDIT.
// Source: job/docker_service.go
//
// Generated by this command:
//
//	mockgen -source job/docker_service.go -destination job/docker_service_mock.go
//

// Package mock_job is a generated GoMock package.
package services

import (
	types "container-manager/types"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDockerService is a mock of DockerService interface.
type MockDockerService struct {
	ctrl     *gomock.Controller
	recorder *MockDockerServiceMockRecorder
}

// MockDockerServiceMockRecorder is the mock recorder for MockDockerService.
type MockDockerServiceMockRecorder struct {
	mock *MockDockerService
}

// NewMockDockerService creates a new mock instance.
func NewMockDockerService(ctrl *gomock.Controller) *MockDockerService {
	mock := &MockDockerService{ctrl: ctrl}
	mock.recorder = &MockDockerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDockerService) EXPECT() *MockDockerServiceMockRecorder {
	return m.recorder
}

// DeployContainer mocks base method.
func (m *MockDockerService) DeployContainer(container types.Container) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployContainer", container)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeployContainer indicates an expected call of DeployContainer.
func (mr *MockDockerServiceMockRecorder) DeployContainer(container any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"DeployContainer",
		reflect.TypeOf((*MockDockerService)(nil).DeployContainer),
		container)
}

// GetContainerStatus mocks base method.
func (m *MockDockerService) GetContainerStatus(containerID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContainerStatus", containerID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContainerStatus indicates an expected call of GetContainerStatus.
func (mr *MockDockerServiceMockRecorder) GetContainerStatus(containerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock,
		"GetContainerStatus",
		reflect.TypeOf((*MockDockerService)(nil).GetContainerStatus),
		containerID)
}
