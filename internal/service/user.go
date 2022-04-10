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
	GetUserData(ctx context.Context, userID string) (*dto.GetUserResponse, error)
	GetUserPosts(ctx context.Context, userID string) (*dto.GetUserPostsResponse, error)
	GetFeed(ctx context.Context, userID string) (*dto.GetUserFeedResponse, error)
	GetProfile(ctx context.Context, request *dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	EditProfile(ctx context.Context, request *dto.EditProfileRequest, userID string) (*dto.EditProfileResponse, error)
}

type userServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *userServiceImpl) GetUserData(ctx context.Context, userID string) (*dto.GetUserResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetUserData success")
	return &dto.GetUserResponse{User: convert.User2DTO(user)}, nil
}

func (svc *userServiceImpl) GetUserPosts(ctx context.Context, userID string) (*dto.GetUserPostsResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	postsCore, err := svc.db.PostRepo.GetPostsByUserID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetPostsByUser error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetUserPosts success")

	posts := []dto.Post{}
	for _, postCore := range postsCore {
		posts = append(posts, convert.Post2DTO(&postCore))
	}

	return &dto.GetUserPostsResponse{Posts: posts}, nil
}

func (svc *userServiceImpl) GetFeed(ctx context.Context, userID string) (*dto.GetUserFeedResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	postsCore, err := svc.db.PostRepo.GetFeed(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetFeed error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetFeed success")

	posts := []dto.Post{}
	for _, postCore := range postsCore {
		posts = append(posts, convert.Post2DTO(&postCore))
	}

	return &dto.GetUserFeedResponse{Posts: posts}, nil
}

func (svc *userServiceImpl) GetProfile(ctx context.Context, request *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	friends, err := svc.db.FriendsRepo.GetFriendsByID(ctx, user.ID)
	if err != nil {
		svc.log.Errorf("GetFriendsByID error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetProfile success")
	return &dto.GetProfileResponse{UserProfile: convert.Profile2DTO(user, friends)}, nil
}

func (svc *userServiceImpl) EditProfile(ctx context.Context, request *dto.EditProfileRequest, userID string) (*dto.EditProfileResponse, error) {
	newUserInfo := convert.EditProfile2Core(&request.NewInfo)

	user, err := svc.db.UserRepo.EditInfo(ctx, newUserInfo, userID)
	if err != nil {
		svc.log.Errorf("EditInfo error: %s", err)
		return nil, err
	}

	svc.log.Debugf("User data after edit: Name: %s; Location: %s; BirthDay: %s; Phone: %s",
		user.Name.Full(), user.Location, user.BirthDay, user.Phone)

	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetFriendsByUserID error: %s", err)
		return nil, err
	}
	svc.log.Debug("EditProfile success")
	return &dto.EditProfileResponse{UserProfile: convert.Profile2DTO(user, friends)}, nil
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
