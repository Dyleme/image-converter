// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Dyleme/image-coverter/internal/service (interfaces: ImageProcesser)

// Package mock_service is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/Dyleme/image-coverter/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockImageProcesser is a mock of ImageProcesser interface.
type MockImageProcesser struct {
	ctrl     *gomock.Controller
	recorder *MockImageProcesserMockRecorder
}

// MockImageProcesserMockRecorder is the mock recorder for MockImageProcesser.
type MockImageProcesserMockRecorder struct {
	mock *MockImageProcesser
}

// NewMockImageProcesser creates a new mock instance.
func NewMockImageProcesser(ctrl *gomock.Controller) *MockImageProcesser {
	mock := &MockImageProcesser{ctrl: ctrl}
	mock.recorder = &MockImageProcesserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageProcesser) EXPECT() *MockImageProcesserMockRecorder {
	return m.recorder
}

// ProcessImage mocks base method.
func (m *MockImageProcesser) ProcessImage(arg0 context.Context, arg1 *model.RequestToProcess) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessImage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessImage indicates an expected call of ProcessImage.
func (mr *MockImageProcesserMockRecorder) ProcessImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessImage", reflect.TypeOf((*MockImageProcesser)(nil).ProcessImage), arg0, arg1)
}
