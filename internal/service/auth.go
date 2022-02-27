//go:generate mockgen -source=user.go -destination=user_mock.go -package=service
package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

type AuthService interface {
	Register(ctx context.Context, request *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
}

type AuthServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *AuthServiceImpl) Register(ctx context.Context, request *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	if exists, err := svc.db.UserRepo.CheckUserEmailExistence(ctx, request.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, constants.ErrEmailAlreadyTaken
	}

	user := &core.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := user.Password.Init(request.Password); err != nil {
		return nil, err
	}

	if err := svc.db.UserRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.RegisterUserResponse{}, nil
}

func NewAuthService(log *logrus.Entry, db *db.Repository) AuthService {
	return &AuthServiceImpl{log: log, db: db}
}
