package auth_service

import (
	auth_db "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/db"
	mock_auth_db "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"testing"
)

// TestRepository ...
type TestRepository struct {
	mockUserR *mock_auth_db.MockAuthRepository
}

// TestRepositories ...
func TestRepositories(t *testing.T, ctrl *gomock.Controller) (*auth_db.Repository, *TestRepository) {
	MockRepo := &TestRepository{
		mock_auth_db.NewMockAuthRepository(ctrl),
	}
	t.Helper()
	return &auth_db.Repository{AuthRepo: MockRepo.mockUserR}, MockRepo
}

// TestLogger ...
func TestLogger(t *testing.T) *logrus.Entry {
	t.Helper()
	logger := logrus.New()
	entry := logrus.NewEntry(logger)
	return entry
}
