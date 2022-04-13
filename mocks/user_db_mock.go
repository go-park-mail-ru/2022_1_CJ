// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/user.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	core "github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// AddDialog mocks base method.
func (m *MockUserRepository) AddDialog(ctx context.Context, dialogID, UserID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDialog", ctx, dialogID, UserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDialog indicates an expected call of AddDialog.
func (mr *MockUserRepositoryMockRecorder) AddDialog(ctx, dialogID, UserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDialog", reflect.TypeOf((*MockUserRepository)(nil).AddDialog), ctx, dialogID, UserID)
}

// CheckUserEmailExistence mocks base method.
func (m *MockUserRepository) CheckUserEmailExistence(ctx context.Context, email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserEmailExistence", ctx, email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserEmailExistence indicates an expected call of CheckUserEmailExistence.
func (mr *MockUserRepositoryMockRecorder) CheckUserEmailExistence(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserEmailExistence", reflect.TypeOf((*MockUserRepository)(nil).CheckUserEmailExistence), ctx, email)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(ctx context.Context, user *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), ctx, user)
}

// DeleteUser mocks base method.
func (m *MockUserRepository) DeleteUser(ctx context.Context, user *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserRepositoryMockRecorder) DeleteUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepository)(nil).DeleteUser), ctx, user)
}

// GetUserByEmail mocks base method.
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepositoryMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(ctx context.Context, ID string) (*core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, ID)
	ret0, _ := ret[0].(*core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), ctx, ID)
}

// GetUserDialogs mocks base method.
func (m *MockUserRepository) GetUserDialogs(ctx context.Context, userID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDialogs", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDialogs indicates an expected call of GetUserDialogs.
func (mr *MockUserRepositoryMockRecorder) GetUserDialogs(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDialogs", reflect.TypeOf((*MockUserRepository)(nil).GetUserDialogs), ctx, userID)
}

// SelectUsers mocks base method.
func (m *MockUserRepository) SelectUsers(ctx context.Context, selector string) ([]core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUsers", ctx, selector)
	ret0, _ := ret[0].([]core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUsers indicates an expected call of SelectUsers.
func (mr *MockUserRepositoryMockRecorder) SelectUsers(ctx, selector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUsers", reflect.TypeOf((*MockUserRepository)(nil).SelectUsers), ctx, selector)
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(ctx context.Context, user *core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), ctx, user)
}

// UserAddPost mocks base method.
func (m *MockUserRepository) UserAddPost(ctx context.Context, userID, postID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserAddPost", ctx, userID, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserAddPost indicates an expected call of UserAddPost.
func (mr *MockUserRepositoryMockRecorder) UserAddPost(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAddPost", reflect.TypeOf((*MockUserRepository)(nil).UserAddPost), ctx, userID, postID)
}

// UserCheckPost mocks base method.
func (m *MockUserRepository) UserCheckPost(ctx context.Context, user *core.User, postID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserCheckPost", ctx, user, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserCheckPost indicates an expected call of UserCheckPost.
func (mr *MockUserRepositoryMockRecorder) UserCheckPost(ctx, user, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserCheckPost", reflect.TypeOf((*MockUserRepository)(nil).UserCheckPost), ctx, user, postID)
}

// UserDeletePost mocks base method.
func (m *MockUserRepository) UserDeletePost(ctx context.Context, userID, postID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserDeletePost", ctx, userID, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserDeletePost indicates an expected call of UserDeletePost.
func (mr *MockUserRepositoryMockRecorder) UserDeletePost(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserDeletePost", reflect.TypeOf((*MockUserRepository)(nil).UserDeletePost), ctx, userID, postID)
}
