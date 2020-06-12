// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prometheus/client_golang/prometheus (interfaces: Counter)

// Package prometheus is a generated GoMock package.
package prometheus

import (
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	reflect "reflect"
)

// MockCounter is a mock of Counter interface
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockCounter) Add(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add
func (mr *MockCounterMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCounter)(nil).Add), arg0)
}

// Collect mocks base method
func (m *MockCounter) Collect(arg0 chan<- prometheus.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Collect", arg0)
}

// Collect indicates an expected call of Collect
func (mr *MockCounterMockRecorder) Collect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockCounter)(nil).Collect), arg0)
}

// Desc mocks base method
func (m *MockCounter) Desc() *prometheus.Desc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Desc")
	ret0, _ := ret[0].(*prometheus.Desc)
	return ret0
}

// Desc indicates an expected call of Desc
func (mr *MockCounterMockRecorder) Desc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Desc", reflect.TypeOf((*MockCounter)(nil).Desc))
}

// Describe mocks base method
func (m *MockCounter) Describe(arg0 chan<- *prometheus.Desc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Describe", arg0)
}

// Describe indicates an expected call of Describe
func (mr *MockCounterMockRecorder) Describe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockCounter)(nil).Describe), arg0)
}

// Inc mocks base method
func (m *MockCounter) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc
func (mr *MockCounterMockRecorder) Inc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
}

// Write mocks base method
func (m *MockCounter) Write(arg0 *io_prometheus_client.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockCounterMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockCounter)(nil).Write), arg0)
}
