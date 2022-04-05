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
	GetUserData(ctx context.Context, UserID string) (*dto.GetUserResponse, error)
	GetUserPosts(ctx context.Context, UserID string) (*dto.GetUserPostsResponse, error)
	GetFeed(ctx context.Context, UserID string) (*dto.GetUserFeedResponse, error)
}

type userServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *userServiceImpl) GetUserData(ctx context.Context, UserID string) (*dto.GetUserResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetUserResponse{User: convert.User2DTO(user)}, nil
}

func (svc *userServiceImpl) GetUserPosts(ctx context.Context, UserID string) (*dto.GetUserPostsResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, UserID)
	if err != nil {
		return nil, err
	}

	posts, err := svc.db.UserRepo.GetPostsByUser(ctx, UserID)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserPostsResponse{PostIDs: posts}, nil
}

func (svc *userServiceImpl) GetFeed(ctx context.Context, UserID string) (*dto.GetUserFeedResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, UserID)
	if err != nil {
		return nil, err
	}

	posts, err := svc.db.PostRepo.GetFeed(ctx, UserID)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserFeedResponse{PostIDs: posts}, nil
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
