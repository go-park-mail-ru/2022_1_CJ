// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/chat.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	core "github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	gomock "github.com/golang/mock/gomock"
)

// MockChatRepository is a mock of ChatRepository interface.
type MockChatRepository struct {
	ctrl     *gomock.Controller
	recorder *MockChatRepositoryMockRecorder
}

// MockChatRepositoryMockRecorder is the mock recorder for MockChatRepository.
type MockChatRepositoryMockRecorder struct {
	mock *MockChatRepository
}

// NewMockChatRepository creates a new mock instance.
func NewMockChatRepository(ctrl *gomock.Controller) *MockChatRepository {
	mock := &MockChatRepository{ctrl: ctrl}
	mock.recorder = &MockChatRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatRepository) EXPECT() *MockChatRepositoryMockRecorder {
	return m.recorder
}

// CreateDialog mocks base method.
func (m *MockChatRepository) CreateDialog(ctx context.Context, userID, name string, authorIDs []string) (*core.Dialog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDialog", ctx, userID, name, authorIDs)
	ret0, _ := ret[0].(*core.Dialog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDialog indicates an expected call of CreateDialog.
func (mr *MockChatRepositoryMockRecorder) CreateDialog(ctx, userID, name, authorIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDialog", reflect.TypeOf((*MockChatRepository)(nil).CreateDialog), ctx, userID, name, authorIDs)
}

// GetDialogByID mocks base method.
func (m *MockChatRepository) GetDialogByID(ctx context.Context, dialogID string) (*core.Dialog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDialogByID", ctx, dialogID)
	ret0, _ := ret[0].(*core.Dialog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDialogByID indicates an expected call of GetDialogByID.
func (mr *MockChatRepositoryMockRecorder) GetDialogByID(ctx, dialogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDialogByID", reflect.TypeOf((*MockChatRepository)(nil).GetDialogByID), ctx, dialogID)
}

// IsChatExist mocks base method.
func (m *MockChatRepository) IsChatExist(ctx context.Context, dialogID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsChatExist", ctx, dialogID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsChatExist indicates an expected call of IsChatExist.
func (mr *MockChatRepositoryMockRecorder) IsChatExist(ctx, dialogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsChatExist", reflect.TypeOf((*MockChatRepository)(nil).IsChatExist), ctx, dialogID)
}

// IsDialogExist mocks base method.
func (m *MockChatRepository) IsDialogExist(ctx context.Context, userID1, userID2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDialogExist", ctx, userID1, userID2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsDialogExist indicates an expected call of IsDialogExist.
func (mr *MockChatRepositoryMockRecorder) IsDialogExist(ctx, userID1, userID2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDialogExist", reflect.TypeOf((*MockChatRepository)(nil).IsDialogExist), ctx, userID1, userID2)
}

// IsUniqDialog mocks base method.
func (m *MockChatRepository) IsUniqDialog(ctx context.Context, userID1, userID2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUniqDialog", ctx, userID1, userID2)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsUniqDialog indicates an expected call of IsUniqDialog.
func (mr *MockChatRepositoryMockRecorder) IsUniqDialog(ctx, userID1, userID2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUniqDialog", reflect.TypeOf((*MockChatRepository)(nil).IsUniqDialog), ctx, userID1, userID2)
}

// ReadMessage mocks base method.
func (m *MockChatRepository) ReadMessage(ctx context.Context, userID, messageID, dialogID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage", ctx, userID, messageID, dialogID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockChatRepositoryMockRecorder) ReadMessage(ctx, userID, messageID, dialogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockChatRepository)(nil).ReadMessage), ctx, userID, messageID, dialogID)
}

// SendMessage mocks base method.
func (m *MockChatRepository) SendMessage(ctx context.Context, message core.Message, dialogID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, message, dialogID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockChatRepositoryMockRecorder) SendMessage(ctx, message, dialogID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockChatRepository)(nil).SendMessage), ctx, message, dialogID)
}
