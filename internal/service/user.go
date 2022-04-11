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
	UpdatePhoto(ctx context.Context, url string, userID string) (*dto.UpdatePhotoResponse, error)
	SearchUsers(ctx context.Context, request *dto.SearchUsersRequest) (*dto.SearchUsersResponse, error)
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
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
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
		posts = append(posts, convert.Post2DTO(&postCore, user))
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
		author, err := svc.db.UserRepo.GetUserByID(ctx, postCore.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, convert.Post2DTO(&postCore, author))
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
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(request.BirthDay) != 0 {
		user.BirthDay = request.BirthDay
	}

	if len(request.Phone) != 0 {
		user.Phone = request.Phone
	}

	if len(request.Location) != 0 {
		user.Location = request.Location
	}

	if len(request.Avatar) != 0 {
		user.Image = request.Avatar
	}

	if len(request.Name.First) != 0 {
		user.Name.First = request.Name.First
	}

	if len(request.Name.Last) != 0 {
		user.Name.Last = request.Name.Last
	}

	if err = svc.db.UserRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.EditProfileResponse{}, nil
}

func (svc *userServiceImpl) UpdatePhoto(ctx context.Context, url string, userID string) (*dto.UpdatePhotoResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Image = url
	if err = svc.db.UserRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UpdatePhotoResponse{URL: url}, nil
}

func (svc *userServiceImpl) SearchUsers(ctx context.Context, request *dto.SearchUsersRequest) (*dto.SearchUsersResponse, error) {
	usersCore, err := svc.db.UserRepo.SelectUsers(ctx, request.Selector)
	if err != nil {
		return nil, err
	}

	users := []dto.User{}
	for _, userCore := range usersCore {
		users = append(users, convert.User2DTO(&userCore))
	}

	return &dto.SearchUsersResponse{Users: users}, nil
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
