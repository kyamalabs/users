// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kyamalabs/users/internal/services (interfaces: AuthGrpcService)
//
// Generated by this command:
//
//	mockgen -package=mockservices -destination=internal/services/mock/auth.go github.com/kyamalabs/users/internal/services AuthGrpcService
//

// Package mockservices is a generated GoMock package.
package mockservices

import (
	reflect "reflect"

	pb "github.com/kyamalabs/auth/api/pb"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthGrpcService is a mock of AuthGrpcService interface.
type MockAuthGrpcService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthGrpcServiceMockRecorder
}

// MockAuthGrpcServiceMockRecorder is the mock recorder for MockAuthGrpcService.
type MockAuthGrpcServiceMockRecorder struct {
	mock *MockAuthGrpcService
}

// NewMockAuthGrpcService creates a new mock instance.
func NewMockAuthGrpcService(ctrl *gomock.Controller) *MockAuthGrpcService {
	mock := &MockAuthGrpcService{ctrl: ctrl}
	mock.recorder = &MockAuthGrpcServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthGrpcService) EXPECT() *MockAuthGrpcServiceMockRecorder {
	return m.recorder
}

// VerifyAccessToken mocks base method.
func (m *MockAuthGrpcService) VerifyAccessToken(arg0 *pb.VerifyAccessTokenRequest, arg1 string) (*pb.VerifyAccessTokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyAccessToken", arg0, arg1)
	ret0, _ := ret[0].(*pb.VerifyAccessTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyAccessToken indicates an expected call of VerifyAccessToken.
func (mr *MockAuthGrpcServiceMockRecorder) VerifyAccessToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyAccessToken", reflect.TypeOf((*MockAuthGrpcService)(nil).VerifyAccessToken), arg0, arg1)
}
