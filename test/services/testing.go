package service

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	mockDB "github.com/go-park-mail-ru/2022_1_CJ/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"testing"
)

// TestRepository ...
type TestRepository struct {
	mockUserR   *mockDB.MockUserRepository
	mockFriensR *mockDB.MockFriendsRepository
	mockPostR   *mockDB.MockPostRepository
	mockChatR   *mockDB.MockChatRepository
}

// TestRepositories ...
func TestRepositories(t *testing.T, ctrl *gomock.Controller) (*db.Repository, *TestRepository) {
	MockRepo := &TestRepository{
		mockDB.NewMockUserRepository(ctrl),
		mockDB.NewMockFriendsRepository(ctrl),
		mockDB.NewMockPostRepository(ctrl),
		mockDB.NewMockChatRepository(ctrl),
	}
	t.Helper()
	return &db.Repository{MockRepo.mockUserR, MockRepo.mockFriensR, MockRepo.mockPostR, MockRepo.mockChatR}, MockRepo
}

// TestLogger ...
func TestLogger(t *testing.T) *logrus.Entry {
	t.Helper()
	logger := logrus.New()
	entry := logrus.NewEntry(logger)
	return entry
}
