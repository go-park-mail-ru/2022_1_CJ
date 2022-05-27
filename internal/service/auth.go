//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
)

type AuthService interface {
	SignupUser(ctx context.Context, request *dto.SignupUserRequest, userID string, token string) (*dto.SignupUserResponse, error)
	LoginUser(ctx context.Context, userID string, token string) (*dto.LoginUserResponse, error)
}

type AuthServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *AuthServiceImpl) SignupUser(ctx context.Context, request *dto.SignupUserRequest, userID string, token string) (*dto.SignupUserResponse, error) {
	user := &core.User{
		ID:    userID,
		Name:  request.Name,
		Email: request.Email,
	}

	if err := svc.db.UserRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	if err := svc.db.FriendsRepo.CreateFriends(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("CreateFriends: %w", err)
	}

	csrfToken, err := utils.GenerateCSRFToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("GenerateCSRFToken: %w", err)
	}

	return &dto.SignupUserResponse{AuthToken: token, CSRFToken: csrfToken}, nil
}

func (svc *AuthServiceImpl) LoginUser(ctx context.Context, userID string, token string) (*dto.LoginUserResponse, error) {
	csrfToken, err := utils.GenerateCSRFToken(userID)
	if err != nil {
		return nil, fmt.Errorf("GenerateCSRFToken: %w", err)
	}
	return &dto.LoginUserResponse{AuthToken: token, CSRFToken: csrfToken}, nil
}

func NewAuthService(log *logrus.Entry, db *db.Repository) AuthService {
	return &AuthServiceImpl{log: log, db: db}
}
