// Code generated by MockGen. DO NOT EDIT.
// Source: gacha.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	service "20dojo-online/pkg/server/service"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGachaServiceInterface is a mock of GachaServiceInterface interface.
type MockGachaServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGachaServiceInterfaceMockRecorder
}

// MockGachaServiceInterfaceMockRecorder is the mock recorder for MockGachaServiceInterface.
type MockGachaServiceInterfaceMockRecorder struct {
	mock *MockGachaServiceInterface
}

// NewMockGachaServiceInterface creates a new mock instance.
func NewMockGachaServiceInterface(ctrl *gomock.Controller) *MockGachaServiceInterface {
	mock := &MockGachaServiceInterface{ctrl: ctrl}
	mock.recorder = &MockGachaServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGachaServiceInterface) EXPECT() *MockGachaServiceInterfaceMockRecorder {
	return m.recorder
}

// DrawGacha mocks base method.
func (m *MockGachaServiceInterface) DrawGacha(serviceRequest *service.DrawGachaRequest) (*service.DrawGachaResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DrawGacha", serviceRequest)
	ret0, _ := ret[0].(*service.DrawGachaResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DrawGacha indicates an expected call of DrawGacha.
func (mr *MockGachaServiceInterfaceMockRecorder) DrawGacha(serviceRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DrawGacha", reflect.TypeOf((*MockGachaServiceInterface)(nil).DrawGacha), serviceRequest)
}
