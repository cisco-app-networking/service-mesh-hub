// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go

// Package mock_zephyr_security is a generated GoMock package.
package mock_zephyr_security

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/solo-io/mesh-projects/pkg/api/security.zephyr.solo.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockVirtualMeshCSRClient is a mock of VirtualMeshCSRClient interface
type MockVirtualMeshCSRClient struct {
	ctrl     *gomock.Controller
	recorder *MockVirtualMeshCSRClientMockRecorder
}

// MockVirtualMeshCSRClientMockRecorder is the mock recorder for MockVirtualMeshCSRClient
type MockVirtualMeshCSRClientMockRecorder struct {
	mock *MockVirtualMeshCSRClient
}

// NewMockVirtualMeshCSRClient creates a new mock instance
func NewMockVirtualMeshCSRClient(ctrl *gomock.Controller) *MockVirtualMeshCSRClient {
	mock := &MockVirtualMeshCSRClient{ctrl: ctrl}
	mock.recorder = &MockVirtualMeshCSRClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVirtualMeshCSRClient) EXPECT() *MockVirtualMeshCSRClientMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockVirtualMeshCSRClient) Create(ctx context.Context, csr *v1alpha1.VirtualMeshCertificateSigningRequest, opts ...client.CreateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, csr}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockVirtualMeshCSRClientMockRecorder) Create(ctx, csr interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, csr}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVirtualMeshCSRClient)(nil).Create), varargs...)
}

// Update mocks base method
func (m *MockVirtualMeshCSRClient) Update(ctx context.Context, csr *v1alpha1.VirtualMeshCertificateSigningRequest, opts ...client.UpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, csr}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockVirtualMeshCSRClientMockRecorder) Update(ctx, csr interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, csr}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVirtualMeshCSRClient)(nil).Update), varargs...)
}

// UpdateStatus mocks base method
func (m *MockVirtualMeshCSRClient) UpdateStatus(ctx context.Context, csr *v1alpha1.VirtualMeshCertificateSigningRequest, opts ...client.UpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, csr}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateStatus", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus
func (mr *MockVirtualMeshCSRClientMockRecorder) UpdateStatus(ctx, csr interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, csr}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockVirtualMeshCSRClient)(nil).UpdateStatus), varargs...)
}

// Get mocks base method
func (m *MockVirtualMeshCSRClient) Get(ctx context.Context, name, namespace string) (*v1alpha1.VirtualMeshCertificateSigningRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, name, namespace)
	ret0, _ := ret[0].(*v1alpha1.VirtualMeshCertificateSigningRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockVirtualMeshCSRClientMockRecorder) Get(ctx, name, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVirtualMeshCSRClient)(nil).Get), ctx, name, namespace)
}

// List mocks base method
func (m *MockVirtualMeshCSRClient) List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.VirtualMeshCertificateSigningRequestList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, opts)
	ret0, _ := ret[0].(*v1alpha1.VirtualMeshCertificateSigningRequestList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockVirtualMeshCSRClientMockRecorder) List(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockVirtualMeshCSRClient)(nil).List), ctx, opts)
}

// Delete mocks base method
func (m *MockVirtualMeshCSRClient) Delete(ctx context.Context, csr *v1alpha1.VirtualMeshCertificateSigningRequest, opts ...client.DeleteOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, csr}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockVirtualMeshCSRClientMockRecorder) Delete(ctx, csr interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, csr}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVirtualMeshCSRClient)(nil).Delete), varargs...)
}
