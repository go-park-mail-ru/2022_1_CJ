//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

type UserService interface {
	GetUserData(ctx context.Context, userID string) (*dto.GetUserResponse, error)
	GetUserPosts(ctx context.Context, request *dto.GetUserPostsRequest) (*dto.GetUserPostsResponse, error)
	GetFeed(ctx context.Context, userID string, request *dto.GetUserFeedRequest) (*dto.GetUserFeedResponse, error)
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

func (svc *userServiceImpl) GetUserPosts(ctx context.Context, request *dto.GetUserPostsRequest) (*dto.GetUserPostsResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	postsCore, pages, err := svc.db.PostRepo.GetPostsByUserID(ctx, request.UserID, request.Page, request.Limit)
	if err != nil {
		svc.log.Errorf("GetPostsByUser error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetUserPosts success")

	var posts []dto.GetPosts
	for _, postCore := range postsCore {
		like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, postCore.ID)
		if err != nil {
			svc.log.Errorf("GetLikeBySubjectID error: %s", err)
			return nil, err
		}

		posts = append(posts, dto.GetPosts{Post: convert.Post2DTOByUser(&postCore, user), Likes: convert.Like2DTO(like, request.UserID)})
	}
	return &dto.GetUserPostsResponse{Posts: posts, Total: pages.Total, AmountPages: pages.AmountPages}, nil
}

func (svc *userServiceImpl) GetFeed(ctx context.Context, userID string, request *dto.GetUserFeedRequest) (*dto.GetUserFeedResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	postsCore, page, err := svc.db.PostRepo.GetFeed(ctx, userID, request.Page, request.Limit)
	if err != nil {
		svc.log.Errorf("GetFeed error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetFeed success")

	var posts []dto.GetPosts
	for _, postCore := range postsCore {
		like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, postCore.ID)
		if err != nil {
			svc.log.Errorf("GetLikeBySubjectID error: %s", err)
			return nil, err
		}
		switch postCore.Type {
		case constants.UserPost:
			author, errUser := svc.db.UserRepo.GetUserByID(ctx, postCore.AuthorID)
			if errUser != nil {
				svc.log.Errorf("GetUserByID error: %s", err)
				return nil, err
			}
			posts = append(posts, dto.GetPosts{Post: convert.Post2DTOByUser(&postCore, author), Likes: convert.Like2DTO(like, userID)})

		case constants.CommunityPost:
			community, errComm := svc.db.CommunityRepo.GetCommunityByID(ctx, postCore.AuthorID)
			if errComm != nil {
				svc.log.Errorf("GetUserByID error: %s", err)
				return nil, err
			}
			var admins []dto.User
			for _, id := range community.AdminIDs {
				user, err := svc.db.UserRepo.GetUserByID(ctx, id)
				if err != nil {
					svc.log.Errorf("GetCommunityByID error: %s", err)
					return nil, constants.ErrDBNotFound
				}
				admins = append(admins, convert.User2DTO(user))
			}
			posts = append(posts, dto.GetPosts{Post: convert.Post2DTOByCommunity(&postCore, community, admins), Likes: convert.Like2DTO(like, userID)})
		default:
			return nil, constants.ErrDBNotFound
		}

	}
	return &dto.GetUserFeedResponse{Posts: posts, Total: page.Total, AmountPages: page.AmountPages}, nil
}

func (svc *userServiceImpl) GetProfile(ctx context.Context, request *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	svc.log.Debug("GetProfile success")
	return &dto.GetProfileResponse{UserProfile: convert.Profile2DTO(user)}, nil
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
		svc.log.Errorf("GetFriendsByUserID error: %s", err)
		return nil, err
	}

	user.Image = url
	if err = svc.db.UserRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UpdatePhotoResponse{URL: url}, nil
}

func (svc *userServiceImpl) SearchUsers(ctx context.Context, request *dto.SearchUsersRequest) (*dto.SearchUsersResponse, error) {
	usersCore, pages, err := svc.db.UserRepo.SelectUsers(ctx, request.Selector, request.Page, request.Limit)
	if err != nil {
		return nil, err
	}

	var users []dto.User
	for _, userCore := range usersCore {
		users = append(users, convert.User2DTO(userCore))
	}
	return &dto.SearchUsersResponse{Users: users, Total: pages.Total, AmountPages: pages.Total}, nil
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
