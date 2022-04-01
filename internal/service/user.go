//go:generate mockgen -source=user.go -destination=user_mock.go -package=service
package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

type UserService interface {
	GetUserData(ctx context.Context, request *dto.GetUserDataRequest) (*dto.GetUserDataResponse, error)
	GetUserFeed(ctx context.Context, request *dto.GetUserFeedRequest) (*dto.GetUserFeedResponse, error)
}

type userServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *userServiceImpl) GetUserData(ctx context.Context, request *dto.GetUserDataRequest) (*dto.GetUserDataResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetUserDataResponse{User: convert.User2DTO(user)}, nil
}

func (svc *userServiceImpl) GetUserFeed(ctx context.Context, request *dto.GetUserFeedRequest) (*dto.GetUserFeedResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	posts, err := svc.db.PostRepo.GetPostsByUser(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetUserFeedResponse{Posts: convert.Posts2DTO(posts)}, nil
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
