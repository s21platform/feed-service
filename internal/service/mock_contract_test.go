// Code generated by MockGen. DO NOT EDIT.
// Source: ./contract.go
//
// Generated by this command:
//
//	mockgen -source=./contract.go -destination=./mock_contract_test.go -package=service
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDBRepo is a mock of DBRepo interface.
type MockDBRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDBRepoMockRecorder
	isgomock struct{}
}

// MockDBRepoMockRecorder is the mock recorder for MockDBRepo.
type MockDBRepoMockRecorder struct {
	mock *MockDBRepo
}

// NewMockDBRepo creates a new mock instance.
func NewMockDBRepo(ctrl *gomock.Controller) *MockDBRepo {
	mock := &MockDBRepo{ctrl: ctrl}
	mock.recorder = &MockDBRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBRepo) EXPECT() *MockDBRepoMockRecorder {
	return m.recorder
}

// Post mocks base method.
func (m *MockDBRepo) Post(ctx context.Context, uuid, content string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", ctx, uuid, content)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Post indicates an expected call of Post.
func (mr *MockDBRepoMockRecorder) Post(ctx, uuid, content any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockDBRepo)(nil).Post), ctx, uuid, content)
}
