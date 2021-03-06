// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Dyleme/image-coverter/pkg/service (interfaces: Autharizater)

// Package mock_service is a generated GoMock package.
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

// GetPasswordHashAndID mocks base method.
func (m *MockAutharizater) GetPasswordHashAndID(arg0 context.Context, arg1 string) ([]byte, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordHashAndID", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPasswordHashAndID indicates an expected call of GetPasswordHashAndID.
func (mr *MockAutharizaterMockRecorder) GetPasswordHashAndID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordHashAndID", reflect.TypeOf((*MockAutharizater)(nil).GetPasswordHashAndID), arg0, arg1)
}
