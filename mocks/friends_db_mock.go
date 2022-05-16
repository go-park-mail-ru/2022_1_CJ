// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/friends.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFriendsRepository is a mock of FriendsRepository interface.
type MockFriendsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFriendsRepositoryMockRecorder
}

// MockFriendsRepositoryMockRecorder is the mock recorder for MockFriendsRepository.
type MockFriendsRepositoryMockRecorder struct {
	mock *MockFriendsRepository
}

// NewMockFriendsRepository creates a new mock instance.
func NewMockFriendsRepository(ctrl *gomock.Controller) *MockFriendsRepository {
	mock := &MockFriendsRepository{ctrl: ctrl}
	mock.recorder = &MockFriendsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriendsRepository) EXPECT() *MockFriendsRepositoryMockRecorder {
	return m.recorder
}

// CreateFriends mocks base method.
func (m *MockFriendsRepository) CreateFriends(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFriends", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFriends indicates an expected call of CreateFriends.
func (mr *MockFriendsRepositoryMockRecorder) CreateFriends(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFriends", reflect.TypeOf((*MockFriendsRepository)(nil).CreateFriends), ctx, userID)
}

// DeleteFriend mocks base method.
func (m *MockFriendsRepository) DeleteFriend(ctx context.Context, exFriendID1, exFriendID2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFriend", ctx, exFriendID1, exFriendID2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFriend indicates an expected call of DeleteFriend.
func (mr *MockFriendsRepositoryMockRecorder) DeleteFriend(ctx, exFriendID1, exFriendID2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFriend", reflect.TypeOf((*MockFriendsRepository)(nil).DeleteFriend), ctx, exFriendID1, exFriendID2)
}

// DeleteIncomingRequest mocks base method.
func (m *MockFriendsRepository) DeleteIncomingRequest(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIncomingRequest", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIncomingRequest indicates an expected call of DeleteIncomingRequest.
func (mr *MockFriendsRepositoryMockRecorder) DeleteIncomingRequest(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIncomingRequest", reflect.TypeOf((*MockFriendsRepository)(nil).DeleteIncomingRequest), ctx, userID, personID)
}

// DeleteOutcomingRequest mocks base method.
func (m *MockFriendsRepository) DeleteOutcomingRequest(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOutcomingRequest", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOutcomingRequest indicates an expected call of DeleteOutcomingRequest.
func (mr *MockFriendsRepositoryMockRecorder) DeleteOutcomingRequest(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOutcomingRequest", reflect.TypeOf((*MockFriendsRepository)(nil).DeleteOutcomingRequest), ctx, userID, personID)
}

// GetFriendsByUserID mocks base method.
func (m *MockFriendsRepository) GetFriendsByUserID(ctx context.Context, userID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendsByUserID", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendsByUserID indicates an expected call of GetFriendsByUserID.
func (mr *MockFriendsRepositoryMockRecorder) GetFriendsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendsByUserID", reflect.TypeOf((*MockFriendsRepository)(nil).GetFriendsByUserID), ctx, userID)
}

// GetIncomingRequestsByUserID mocks base method.
func (m *MockFriendsRepository) GetIncomingRequestsByUserID(ctx context.Context, userID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIncomingRequestsByUserID", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIncomingRequestsByUserID indicates an expected call of GetIncomingRequestsByUserID.
func (mr *MockFriendsRepositoryMockRecorder) GetIncomingRequestsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIncomingRequestsByUserID", reflect.TypeOf((*MockFriendsRepository)(nil).GetIncomingRequestsByUserID), ctx, userID)
}

// GetOutcomingRequestsByUserID mocks base method.
func (m *MockFriendsRepository) GetOutcomingRequestsByUserID(ctx context.Context, userID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutcomingRequestsByUserID", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutcomingRequestsByUserID indicates an expected call of GetOutcomingRequestsByUserID.
func (mr *MockFriendsRepositoryMockRecorder) GetOutcomingRequestsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutcomingRequestsByUserID", reflect.TypeOf((*MockFriendsRepository)(nil).GetOutcomingRequestsByUserID), ctx, userID)
}

// IsNotFriend mocks base method.
func (m *MockFriendsRepository) IsNotFriend(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNotFriend", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsNotFriend indicates an expected call of IsNotFriend.
func (mr *MockFriendsRepositoryMockRecorder) IsNotFriend(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNotFriend", reflect.TypeOf((*MockFriendsRepository)(nil).IsNotFriend), ctx, userID, personID)
}

// IsUniqRequest mocks base method.
func (m *MockFriendsRepository) IsUniqRequest(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUniqRequest", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsUniqRequest indicates an expected call of IsUniqRequest.
func (mr *MockFriendsRepositoryMockRecorder) IsUniqRequest(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUniqRequest", reflect.TypeOf((*MockFriendsRepository)(nil).IsUniqRequest), ctx, userID, personID)
}

// MakeFriends mocks base method.
func (m *MockFriendsRepository) MakeFriends(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeFriends", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeFriends indicates an expected call of MakeFriends.
func (mr *MockFriendsRepositoryMockRecorder) MakeFriends(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeFriends", reflect.TypeOf((*MockFriendsRepository)(nil).MakeFriends), ctx, userID, personID)
}

// MakeIncomingRequest mocks base method.
func (m *MockFriendsRepository) MakeIncomingRequest(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeIncomingRequest", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeIncomingRequest indicates an expected call of MakeIncomingRequest.
func (mr *MockFriendsRepositoryMockRecorder) MakeIncomingRequest(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeIncomingRequest", reflect.TypeOf((*MockFriendsRepository)(nil).MakeIncomingRequest), ctx, userID, personID)
}

// MakeOutcomingRequest mocks base method.
func (m *MockFriendsRepository) MakeOutcomingRequest(ctx context.Context, userID, personID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeOutcomingRequest", ctx, userID, personID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeOutcomingRequest indicates an expected call of MakeOutcomingRequest.
func (mr *MockFriendsRepositoryMockRecorder) MakeOutcomingRequest(ctx, userID, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeOutcomingRequest", reflect.TypeOf((*MockFriendsRepository)(nil).MakeOutcomingRequest), ctx, userID, personID)
}
