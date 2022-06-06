// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/lang/frontend/starlark/interpreter.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInterpreter is a mock of Interpreter interface.
type MockInterpreter struct {
	ctrl     *gomock.Controller
	recorder *MockInterpreterMockRecorder
}

// MockInterpreterMockRecorder is the mock recorder for MockInterpreter.
type MockInterpreterMockRecorder struct {
	mock *MockInterpreter
}

// NewMockInterpreter creates a new mock instance.
func NewMockInterpreter(ctrl *gomock.Controller) *MockInterpreter {
	mock := &MockInterpreter{ctrl: ctrl}
	mock.recorder = &MockInterpreterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterpreter) EXPECT() *MockInterpreterMockRecorder {
	return m.recorder
}

// Eval mocks base method.
func (m *MockInterpreter) Eval(script string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Eval", script)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Eval indicates an expected call of Eval.
func (mr *MockInterpreterMockRecorder) Eval(script interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Eval", reflect.TypeOf((*MockInterpreter)(nil).Eval), script)
}

// ExecFile mocks base method.
func (m *MockInterpreter) ExecFile(filename, funcname string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecFile", filename, funcname)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecFile indicates an expected call of ExecFile.
func (mr *MockInterpreterMockRecorder) ExecFile(filename, funcname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecFile", reflect.TypeOf((*MockInterpreter)(nil).ExecFile), filename, funcname)
}
