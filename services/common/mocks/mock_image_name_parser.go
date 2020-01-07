// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/solo-io/mesh-projects/services/common (interfaces: ImageNameParser)

// Package mock_common is a generated GoMock package.
package mock_common

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/solo-io/mesh-projects/services/common"
)

// MockImageNameParser is a mock of ImageNameParser interface
type MockImageNameParser struct {
	ctrl     *gomock.Controller
	recorder *MockImageNameParserMockRecorder
}

// MockImageNameParserMockRecorder is the mock recorder for MockImageNameParser
type MockImageNameParserMockRecorder struct {
	mock *MockImageNameParser
}

// NewMockImageNameParser creates a new mock instance
func NewMockImageNameParser(ctrl *gomock.Controller) *MockImageNameParser {
	mock := &MockImageNameParser{ctrl: ctrl}
	mock.recorder = &MockImageNameParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImageNameParser) EXPECT() *MockImageNameParserMockRecorder {
	return m.recorder
}

// Parse mocks base method
func (m *MockImageNameParser) Parse(arg0 string) (*common.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0)
	ret0, _ := ret[0].(*common.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse
func (mr *MockImageNameParserMockRecorder) Parse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockImageNameParser)(nil).Parse), arg0)
}