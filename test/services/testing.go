package service

import (
	"testing"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	mockDB "github.com/go-park-mail-ru/2022_1_CJ/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

// TestRepository ...
type TestRepository struct {
	mockUserR    *mockDB.MockUserRepository
	mockFriendsR *mockDB.MockFriendsRepository
	mockPostR    *mockDB.MockPostRepository
	mockChatR    *mockDB.MockChatRepository
	mockLikeR    *mockDB.MockLikeRepository
	// mockCommunityR *mockDB.MockCommunityRepository
}

// TestRepositories ...
func TestRepositories(t *testing.T, ctrl *gomock.Controller) (*db.Repository, *TestRepository) {
	MockRepo := &TestRepository{
		mockDB.NewMockUserRepository(ctrl),
		mockDB.NewMockFriendsRepository(ctrl),
		mockDB.NewMockPostRepository(ctrl),
		mockDB.NewMockChatRepository(ctrl),
		mockDB.NewMockLikeRepository(ctrl),
		// mockDB.NewMockCommunityRepository(ctrl),
	}
	t.Helper()
	return &db.Repository{UserRepo: MockRepo.mockUserR,
		FriendsRepo: MockRepo.mockFriendsR,
		PostRepo:    MockRepo.mockPostR,
		ChatRepo:    MockRepo.mockChatR,
		LikeRepo:    MockRepo.mockLikeR,
		// CommunityRepo: MockRepo.mockCommunityR
	}, MockRepo
}

// TestLogger ...
func TestLogger(t *testing.T) *logrus.Entry {
	t.Helper()
	logger := logrus.New()
	entry := logrus.NewEntry(logger)
	return entry
}
