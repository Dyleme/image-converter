// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Dyleme/image-coverter/internal/handler (interfaces: Autharizater)

// Package mock_handler is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/Dyleme/image-coverter/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockAutharizater is a mock of Autharizater interface.
type MockAutharizater struct {
	ctrl     *gomock.Controller
	recorder *MockAutharizaterMockRecorder
}

// MockAutharizaterMockRecorder is the mock recorder for MockAutharizater.
type MockAutharizaterMockRecorder struct {
	mock *MockAutharizater
}

// NewMockAutharizater creates a new mock instance.
func NewMockAutharizater(ctrl *gomock.Controller) *MockAutharizater {
	mock := &MockAutharizater{ctrl: ctrl}
	mock.recorder = &MockAutharizaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutharizater) EXPECT() *MockAutharizaterMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAutharizater) CreateUser(arg0 context.Context, arg1 model.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAutharizaterMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAutharizater)(nil).CreateUser), arg0, arg1)
}

// ValidateUser mocks base method.
func (m *MockAutharizater) ValidateUser(arg0 context.Context, arg1 model.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUser", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUser indicates an expected call of ValidateUser.
func (mr *MockAutharizaterMockRecorder) ValidateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUser", reflect.TypeOf((*MockAutharizater)(nil).ValidateUser), arg0, arg1)
}