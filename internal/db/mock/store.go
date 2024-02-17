// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kyamalabs/users/internal/db/sqlc (interfaces: Store)
//
// Generated by this command:
//
//	mockgen -package=mockdb -destination=internal/db/mock/store.go github.com/kyamalabs/users/internal/db/sqlc Store
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/kyamalabs/users/internal/db/sqlc"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateProfile mocks base method.
func (m *MockStore) CreateProfile(arg0 context.Context, arg1 db.CreateProfileParams) (db.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfile", arg0, arg1)
	ret0, _ := ret[0].(db.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfile indicates an expected call of CreateProfile.
func (mr *MockStoreMockRecorder) CreateProfile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfile", reflect.TypeOf((*MockStore)(nil).CreateProfile), arg0, arg1)
}

// CreateProfileTx mocks base method.
func (m *MockStore) CreateProfileTx(arg0 context.Context, arg1 db.CreateProfileTxParams) (db.CreateProfileTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfileTx", arg0, arg1)
	ret0, _ := ret[0].(db.CreateProfileTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfileTx indicates an expected call of CreateProfileTx.
func (mr *MockStoreMockRecorder) CreateProfileTx(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfileTx", reflect.TypeOf((*MockStore)(nil).CreateProfileTx), arg0, arg1)
}

// CreateReferral mocks base method.
func (m *MockStore) CreateReferral(arg0 context.Context, arg1 db.CreateReferralParams) (db.Referral, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReferral", arg0, arg1)
	ret0, _ := ret[0].(db.Referral)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReferral indicates an expected call of CreateReferral.
func (mr *MockStoreMockRecorder) CreateReferral(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReferral", reflect.TypeOf((*MockStore)(nil).CreateReferral), arg0, arg1)
}

// DeleteProfile mocks base method.
func (m *MockStore) DeleteProfile(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProfile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProfile indicates an expected call of DeleteProfile.
func (mr *MockStoreMockRecorder) DeleteProfile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProfile", reflect.TypeOf((*MockStore)(nil).DeleteProfile), arg0, arg1)
}

// GetProfile mocks base method.
func (m *MockStore) GetProfile(arg0 context.Context, arg1 string) (db.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", arg0, arg1)
	ret0, _ := ret[0].(db.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockStoreMockRecorder) GetProfile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockStore)(nil).GetProfile), arg0, arg1)
}

// GetProfileTx mocks base method.
func (m *MockStore) GetProfileTx(arg0 context.Context, arg1 db.GetProfileTxParams) (db.GetProfileTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileTx", arg0, arg1)
	ret0, _ := ret[0].(db.GetProfileTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileTx indicates an expected call of GetProfileTx.
func (mr *MockStoreMockRecorder) GetProfileTx(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileTx", reflect.TypeOf((*MockStore)(nil).GetProfileTx), arg0, arg1)
}

// GetProfilesCount mocks base method.
func (m *MockStore) GetProfilesCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfilesCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfilesCount indicates an expected call of GetProfilesCount.
func (mr *MockStoreMockRecorder) GetProfilesCount(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfilesCount", reflect.TypeOf((*MockStore)(nil).GetProfilesCount), arg0)
}

// GetReferrer mocks base method.
func (m *MockStore) GetReferrer(arg0 context.Context, arg1 string) (db.Referral, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReferrer", arg0, arg1)
	ret0, _ := ret[0].(db.Referral)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReferrer indicates an expected call of GetReferrer.
func (mr *MockStoreMockRecorder) GetReferrer(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReferrer", reflect.TypeOf((*MockStore)(nil).GetReferrer), arg0, arg1)
}

// ListProfiles mocks base method.
func (m *MockStore) ListProfiles(arg0 context.Context, arg1 db.ListProfilesParams) ([]db.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProfiles", arg0, arg1)
	ret0, _ := ret[0].([]db.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProfiles indicates an expected call of ListProfiles.
func (mr *MockStoreMockRecorder) ListProfiles(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProfiles", reflect.TypeOf((*MockStore)(nil).ListProfiles), arg0, arg1)
}

// ListReferrals mocks base method.
func (m *MockStore) ListReferrals(arg0 context.Context, arg1 db.ListReferralsParams) ([]db.Referral, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListReferrals", arg0, arg1)
	ret0, _ := ret[0].([]db.Referral)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListReferrals indicates an expected call of ListReferrals.
func (mr *MockStoreMockRecorder) ListReferrals(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListReferrals", reflect.TypeOf((*MockStore)(nil).ListReferrals), arg0, arg1)
}

// UpdateProfile mocks base method.
func (m *MockStore) UpdateProfile(arg0 context.Context, arg1 db.UpdateProfileParams) (db.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", arg0, arg1)
	ret0, _ := ret[0].(db.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockStoreMockRecorder) UpdateProfile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockStore)(nil).UpdateProfile), arg0, arg1)
}

// UpdateProfileTx mocks base method.
func (m *MockStore) UpdateProfileTx(arg0 context.Context, arg1 db.UpdateProfileTxParams) (db.UpdateProfileTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfileTx", arg0, arg1)
	ret0, _ := ret[0].(db.UpdateProfileTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfileTx indicates an expected call of UpdateProfileTx.
func (mr *MockStoreMockRecorder) UpdateProfileTx(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileTx", reflect.TypeOf((*MockStore)(nil).UpdateProfileTx), arg0, arg1)
}
