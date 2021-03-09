// Code generated by MockGen. DO NOT EDIT.
// Source: response.go

// Package mock_response is a generated GoMock package.
package mock_response

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHttpResponseInterface is a mock of HttpResponseInterface interface.
type MockHttpResponseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHttpResponseInterfaceMockRecorder
}

// MockHttpResponseInterfaceMockRecorder is the mock recorder for MockHttpResponseInterface.
type MockHttpResponseInterfaceMockRecorder struct {
	mock *MockHttpResponseInterface
}

// NewMockHttpResponseInterface creates a new mock instance.
func NewMockHttpResponseInterface(ctrl *gomock.Controller) *MockHttpResponseInterface {
	mock := &MockHttpResponseInterface{ctrl: ctrl}
	mock.recorder = &MockHttpResponseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpResponseInterface) EXPECT() *MockHttpResponseInterfaceMockRecorder {
	return m.recorder
}

// Failed mocks base method.
func (m *MockHttpResponseInterface) Failed(writer http.ResponseWriter, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Failed", writer, err)
}

// Failed indicates an expected call of Failed.
func (mr *MockHttpResponseInterfaceMockRecorder) Failed(writer, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Failed", reflect.TypeOf((*MockHttpResponseInterface)(nil).Failed), writer, err)
}

// Success mocks base method.
func (m *MockHttpResponseInterface) Success(writer http.ResponseWriter, response interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Success", writer, response)
}

// Success indicates an expected call of Success.
func (mr *MockHttpResponseInterfaceMockRecorder) Success(writer, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Success", reflect.TypeOf((*MockHttpResponseInterface)(nil).Success), writer, response)
}
