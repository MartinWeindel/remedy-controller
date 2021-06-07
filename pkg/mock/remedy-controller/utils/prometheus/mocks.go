// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/remedy-controller/pkg/utils/prometheus (interfaces: GaugeVec)

// Package prometheus is a generated GoMock package.
package prometheus

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
)

// MockGaugeVec is a mock of GaugeVec interface.
type MockGaugeVec struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeVecMockRecorder
}

// MockGaugeVecMockRecorder is the mock recorder for MockGaugeVec.
type MockGaugeVecMockRecorder struct {
	mock *MockGaugeVec
}

// NewMockGaugeVec creates a new mock instance.
func NewMockGaugeVec(ctrl *gomock.Controller) *MockGaugeVec {
	mock := &MockGaugeVec{ctrl: ctrl}
	mock.recorder = &MockGaugeVecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGaugeVec) EXPECT() *MockGaugeVecMockRecorder {
	return m.recorder
}

// DeleteLabelValues mocks base method.
func (m *MockGaugeVec) DeleteLabelValues(arg0 ...string) bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteLabelValues", varargs...)
	ret0, _ := ret[0].(bool)
	return ret0
}

// DeleteLabelValues indicates an expected call of DeleteLabelValues.
func (mr *MockGaugeVecMockRecorder) DeleteLabelValues(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLabelValues", reflect.TypeOf((*MockGaugeVec)(nil).DeleteLabelValues), arg0...)
}

// WithLabelValues mocks base method.
func (m *MockGaugeVec) WithLabelValues(arg0 ...string) prometheus.Gauge {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithLabelValues", varargs...)
	ret0, _ := ret[0].(prometheus.Gauge)
	return ret0
}

// WithLabelValues indicates an expected call of WithLabelValues.
func (mr *MockGaugeVecMockRecorder) WithLabelValues(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLabelValues", reflect.TypeOf((*MockGaugeVec)(nil).WithLabelValues), arg0...)
}
