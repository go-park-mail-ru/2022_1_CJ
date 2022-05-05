package service

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/sirupsen/logrus"
)

type CommunityService interface {
	CreateCommunity(ctx context.Context, request *dto.CreateCommunityRequest, userID string) (*dto.CreateCommunityResponse, error)
	DeleteCommunity(ctx context.Context, request *dto.DeleteCommunityRequest, userID string) (*dto.DeleteCommunityResponse, error)
	EditCommunity(ctx context.Context, request *dto.EditCommunityRequest, userID string) (*dto.EditCommunityResponse, error)
	GetCommunity(ctx context.Context, request *dto.GetCommunityRequest) (*dto.GetCommunityResponse, error)
	GetCommunityPosts(ctx context.Context, request *dto.GetCommunityPostsRequest, userID string) (*dto.GetCommunityPostsResponse, error)
	GetUserCommunities(ctx context.Context, request *dto.GetUserCommunitiesRequest) (*dto.GetUserCommunitiesResponse, error)
	GetUserManageCommunities(ctx context.Context, request *dto.GetUserManageCommunitiesRequest) (*dto.GetUserManageCommunitiesResponse, error)
	JoinCommunity(ctx context.Context, request *dto.JoinCommunityRequest, userID string) (*dto.JoinCommunityResponse, error)
	LeaveCommunity(ctx context.Context, request *dto.LeaveCommunityRequest, userID string) (*dto.LeaveCommunityResponse, error)
	SearchCommunities(ctx context.Context, request *dto.SearchCommunitiesRequest) (*dto.SearchCommunitiesResponse, error)
	GetFollowers(ctx context.Context, request *dto.GetFollowersRequest) (*dto.GetFollowersResponse, error)
	GetCommunities(ctx context.Context, request *dto.GetCommunitiesRequest) (*dto.GetCommunitiesResponse, error)
	UpdatePhoto(ctx context.Context, request *dto.UpdatePhotoCommunityRequest, url string, userID string) (*dto.UpdatePhotoCommunityResponse, error)
	GetMutualFriends(ctx context.Context, request *dto.GetMutualFriendsRequest, userID string) (*dto.GetMutualFriendsResponse, error)

	CreatePostCommunity(ctx context.Context, request *dto.CreatePostCommunityRequest, userID string) (*dto.CreatePostCommunityResponse, error)
	DeletePostCommunity(ctx context.Context, request *dto.DeletePostCommunityRequest, userID string) (*dto.DeletePostCommunityResponse, error)
	EditPostCommunity(ctx context.Context, request *dto.EditPostCommunityRequest, userID string) (*dto.EditPostCommunityResponse, error)
}

type communityServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *communityServiceImpl) CreateCommunity(ctx context.Context, request *dto.CreateCommunityRequest, userID string) (*dto.CreateCommunityResponse, error) {
	request.Admins = append(request.Admins, userID)
	community, err := svc.db.CommunityRepo.CreateCommunity(ctx, &core.Community{
		Name:        request.Name,
		Image:       request.Image,
		Info:        request.Info,
		AdminIDs:    request.Admins,
		FollowerIDs: request.Admins,
	})
	if err != nil {
		svc.log.Errorf("CreateCommunity error: %s", err)
		return nil, err
	}
	for _, id := range request.Admins {
		err = svc.db.UserRepo.UserAddCommunity(ctx, id, community.ID)
		if err != nil {
			svc.log.Errorf("UserAddCommunity error: %s", err)
			return nil, err
		}
	}
	return &dto.CreateCommunityResponse{}, nil
}

func (svc *communityServiceImpl) GetCommunity(ctx context.Context, request *dto.GetCommunityRequest) (*dto.GetCommunityResponse, error) {
	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
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

	return &dto.GetCommunityResponse{Community: convert.Community2DTOprofile(community, admins)}, nil
}

func (svc *communityServiceImpl) GetUserManageCommunities(ctx context.Context, request *dto.GetUserManageCommunitiesRequest) (*dto.GetUserManageCommunitiesResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	communityIDs, total, pages := utils.GetLimitArray(&user.CommunityIDs, request.Limit, request.Page)
	var communities []dto.Community
	for _, commID := range communityIDs {
		comm, err := svc.db.CommunityRepo.GetCommunityByID(ctx, commID)
		if err != nil {
			svc.log.Errorf("GetCommunityByID error: %s", err)
			return nil, constants.ErrDBNotFound
		}

		for _, id := range comm.AdminIDs {
			if id == request.UserID {
				communities = append(communities, convert.Community2DTO(comm))
			}
		}
	}

	return &dto.GetUserManageCommunitiesResponse{Communities: communities, Total: total, AmountPages: pages}, nil
}

func (svc *communityServiceImpl) GetCommunities(ctx context.Context, request *dto.GetCommunitiesRequest) (*dto.GetCommunitiesResponse, error) {
	communities, page, err := svc.db.CommunityRepo.GetAllCommunities(ctx, request.Limit, request.Page)
	if err != nil {
		svc.log.Errorf("GetAllCommunities error: %s", err)
		return nil, constants.ErrDBNotFound
	}
	var res []dto.Community
	for _, comm := range communities {
		res = append(res, convert.Community2DTO(&comm))
	}
	return &dto.GetCommunitiesResponse{Communities: res, Total: page.Total, AmountPages: page.AmountPages}, nil
}

func (svc *communityServiceImpl) JoinCommunity(ctx context.Context, request *dto.JoinCommunityRequest, userID string) (*dto.JoinCommunityResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	for _, id := range user.CommunityIDs {
		if id == request.CommunityID {
			return nil, constants.ErrAlreadyFollower
		}
	}

	err = svc.db.CommunityRepo.AddFollower(ctx, request.CommunityID, userID)
	if err != nil {
		svc.log.Errorf("AddFollower error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	err = svc.db.UserRepo.UserAddCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserAddCommunity error: %s", err)
		return nil, err
	}

	return &dto.JoinCommunityResponse{}, nil
}

func (svc *communityServiceImpl) LeaveCommunity(ctx context.Context, request *dto.LeaveCommunityRequest, userID string) (*dto.LeaveCommunityResponse, error) {
	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("DeleteFollower error: %s", err)
		return nil, err
	}

	actualList := len(community.AdminIDs)
	for _, id := range community.AdminIDs {
		if id == userID {
			actualList -= 1
			err = svc.db.CommunityRepo.DeleteAdmin(ctx, request.CommunityID, id)
			if err != nil {
				svc.log.Errorf("DeleteAdmin error: %s", err)
				return nil, err
			}
		}
	}

	err = svc.db.CommunityRepo.DeleteFollower(ctx, request.CommunityID, userID)
	if err != nil {
		svc.log.Errorf("DeleteFollower error: %s", err)
		return nil, err
	}

	err = svc.db.UserRepo.UserDeleteCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserDeleteCommunity error: %s", err)
		return nil, err
	}

	if actualList == 0 {
		_, err := svc.DeleteCommunity(ctx, &dto.DeleteCommunityRequest{CommunityID: request.CommunityID}, userID)
		if err != nil {
			svc.log.Errorf("DeleteCommunity error: %s", err)
			return nil, err
		}
	}

	return &dto.LeaveCommunityResponse{}, nil
}

func (svc *communityServiceImpl) GetFollowers(ctx context.Context, request *dto.GetFollowersRequest) (*dto.GetFollowersResponse, error) {
	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	followerIDs, total, pages := utils.GetLimitArray(&community.FollowerIDs, request.Limit, request.Page)
	var followers []dto.User
	for _, id := range followerIDs {
		user, err := svc.db.UserRepo.GetUserByID(ctx, id)
		if err != nil {
			svc.log.Errorf("GetCommunityByID error: %s", err)
			return nil, constants.ErrDBNotFound
		}
		followers = append(followers, convert.User2DTO(user))
	}

	return &dto.GetFollowersResponse{Amount: int64(len(community.FollowerIDs)), Followers: followers, AmountPages: pages, Total: total}, nil
}

func (svc *communityServiceImpl) GetMutualFriends(ctx context.Context, request *dto.GetMutualFriendsRequest, userID string) (*dto.GetMutualFriendsResponse, error) {
	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, user.FriendsID)
	if err != nil {
		svc.log.Errorf("GetFriendsByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	friendsIDs, total, pages := utils.GetLimitArray(&friends, request.Limit, request.Page)

	var followers []dto.User
	for _, id1 := range friendsIDs {
		for _, id2 := range community.FollowerIDs {
			if id1 == id2 {
				userFriend, err := svc.db.UserRepo.GetUserByID(ctx, id1)
				if err != nil {
					svc.log.Errorf("GetCommunityByID error: %s", err)
					return nil, constants.ErrDBNotFound
				}
				followers = append(followers, convert.User2DTO(userFriend))
			}
		}
	}

	return &dto.GetMutualFriendsResponse{Amount: int64(len(followers)), Followers: followers, Total: total, AmountPages: pages}, nil
}

func (svc *communityServiceImpl) GetUserCommunities(ctx context.Context, request *dto.GetUserCommunitiesRequest) (*dto.GetUserCommunitiesResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	communityIDs, total, pages := utils.GetLimitArray(&user.CommunityIDs, request.Limit, request.Page)
	var communities []dto.Community
	for _, id := range communityIDs {
		comm, err := svc.db.CommunityRepo.GetCommunityByID(ctx, id)
		if err != nil {
			svc.log.Errorf("GetCommunityByID error: %s", err)
			return nil, constants.ErrDBNotFound
		}
		communities = append(communities, convert.Community2DTO(comm))
	}
	return &dto.GetUserCommunitiesResponse{Communities: communities, Total: total, AmountPages: pages}, nil
}

func (svc *communityServiceImpl) SearchCommunities(ctx context.Context, request *dto.SearchCommunitiesRequest) (*dto.SearchCommunitiesResponse, error) {
	communities, page, err := svc.db.CommunityRepo.SearchCommunities(ctx, request.Selector, request.Limit, request.Page)
	if err != nil {
		svc.log.Errorf("SearchCommunities error: %s", err)
		return nil, constants.ErrDBNotFound
	}
	var res []dto.Community
	for _, comm := range communities {
		res = append(res, convert.Community2DTO(&comm))
	}
	return &dto.SearchCommunitiesResponse{Communities: res, AmountPages: page.AmountPages, Total: page.Total}, nil
}

func (svc *communityServiceImpl) UpdatePhoto(ctx context.Context, request *dto.UpdatePhotoCommunityRequest, url string, userID string) (*dto.UpdatePhotoCommunityResponse, error) {
	if err := svc.db.UserRepo.UserCheckCommunity(ctx, userID, request.CommunityID); err != nil {
		return nil, fmt.Errorf("UserCheckCommunity: %w", err)
	}

	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		return nil, fmt.Errorf("GetCommunityByID: %w", err)
	}

	for num, id := range community.AdminIDs {
		if id == userID {
			break
		} else if num == (len(community.AdminIDs) - 1) {
			return nil, constants.ErrAuthorIDMismatch
		}
	}

	community.Image = url
	if err = svc.db.CommunityRepo.EditCommunity(ctx, community); err != nil {
		return nil, err
	}

	return &dto.UpdatePhotoCommunityResponse{URL: url}, nil
}

func (svc *communityServiceImpl) GetCommunityPosts(ctx context.Context, request *dto.GetCommunityPostsRequest, userID string) (*dto.GetCommunityPostsResponse, error) {
	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}
	postIDs, total, pages := utils.GetLimitArray(&community.PostIDs, request.Limit, request.Page)
	var posts []dto.GetPosts
	for _, id := range postIDs {
		post, err := svc.db.PostRepo.GetPostByID(ctx, id)
		if err != nil {
			svc.log.Errorf("GetPostByID error: %s", err)
			return nil, constants.ErrDBNotFound
		}
		like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, post.ID)
		if err != nil {
			svc.log.Errorf("GetLikeBySubjectID error: %s", err)
			return nil, err
		}
		posts = append(posts, dto.GetPosts{Post: convert.Post2DTOByCommunity(post, community), Likes: convert.Like2DTO(like, userID)})
	}

	return &dto.GetCommunityPostsResponse{Posts: posts, Total: total, AmountPages: pages}, nil
}

func (svc *communityServiceImpl) EditCommunity(ctx context.Context, request *dto.EditCommunityRequest, userID string) (*dto.EditCommunityResponse, error) {
	request.Admins = append(request.Admins, userID)
	err := svc.db.UserRepo.UserCheckCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserCheckCommunity error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	for num, id := range community.AdminIDs {
		if id == userID {
			break
		} else if num == (len(community.AdminIDs) - 1) {
			return nil, constants.ErrAuthorIDMismatch
		}
	}

	community.Name = request.Name
	community.Image = request.Image
	community.Info = request.Info
	community.AdminIDs = request.Admins

	err = svc.db.CommunityRepo.EditCommunity(ctx, community)
	if err != nil {
		svc.log.Errorf("EditCommunity error: %s", err)
		return nil, err
	}

	return &dto.EditCommunityResponse{}, nil
}

func (svc *communityServiceImpl) CreatePostCommunity(ctx context.Context, request *dto.CreatePostCommunityRequest, userID string) (*dto.CreatePostCommunityResponse, error) {
	err := svc.db.UserRepo.UserCheckCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserCheckCommunity error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	for num, id := range community.AdminIDs {
		if id == userID {
			break
		} else if num == (len(community.AdminIDs) - 1) {
			return nil, constants.ErrAuthorIDMismatch
		}
	}
	post, err := svc.db.PostRepo.CreatePost(ctx, &core.Post{
		AuthorID: community.ID,
		Message:  request.Message,
		Images:   request.Images,
		Type:     constants.CommunityPost,
	})
	if err != nil {
		svc.log.Errorf("CreatePost error: %s", err)
		return nil, err
	}

	err = svc.db.CommunityRepo.CommunityAddPost(ctx, community.ID, post.ID)
	if err != nil {
		svc.log.Errorf("UserAddPost error: %s", err)
		return nil, err
	}

	_, err = svc.db.LikeRepo.CreateLike(ctx, &core.Like{Subject: post.ID})
	if err != nil {
		svc.log.Errorf("CreateLike error: %s", err)
		return nil, err
	}

	return &dto.CreatePostCommunityResponse{}, nil
}

func (svc *communityServiceImpl) EditPostCommunity(ctx context.Context, request *dto.EditPostCommunityRequest, userID string) (*dto.EditPostCommunityResponse, error) {
	err := svc.db.UserRepo.UserCheckCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserCheckCommunity error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	for num, id := range community.AdminIDs {
		if id == userID {
			break
		} else if num == (len(community.AdminIDs) - 1) {
			return nil, constants.ErrAuthorIDMismatch
		}
	}

	_, err = svc.db.PostRepo.EditPost(ctx, &core.Post{
		AuthorID: request.CommunityID,
		ID:       request.PostID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		svc.log.Errorf("EditPost error: %s", err)
		return nil, err
	}

	return &dto.EditPostCommunityResponse{}, nil
}

func (svc *communityServiceImpl) DeletePostCommunity(ctx context.Context, request *dto.DeletePostCommunityRequest, userID string) (*dto.DeletePostCommunityResponse, error) {
	err := svc.db.UserRepo.UserCheckCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserCheckCommunity error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	for num, id := range community.AdminIDs {
		if id == userID {
			break
		} else if num == (len(community.AdminIDs) - 1) {
			return nil, constants.ErrAuthorIDMismatch
		}
	}

	err = svc.db.PostRepo.DeletePost(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("DeletePost error: %s", err)
		return nil, err
	}

	err = svc.db.CommunityRepo.CommunityDeletePost(ctx, community.ID, request.PostID)
	if err != nil {
		svc.log.Errorf("CommunityDeletePost error: %s", err)
		return nil, err
	}

	err = svc.db.LikeRepo.DeleteLike(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("DeleteLike error: %s", err)
		return nil, err
	}

	return &dto.DeletePostCommunityResponse{}, nil
}

func (svc *communityServiceImpl) DeleteCommunity(ctx context.Context, request *dto.DeleteCommunityRequest, userID string) (*dto.DeleteCommunityResponse, error) {
	err := svc.db.UserRepo.UserCheckCommunity(ctx, userID, request.CommunityID)
	if err != nil {
		svc.log.Errorf("UserCheckCommunity error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	for num, id := range community.AdminIDs {
		if id == userID {
			break
		} else if num == (len(community.AdminIDs) - 1) {
			return nil, constants.ErrAuthorIDMismatch
		}
	}

	err = svc.db.CommunityRepo.DeleteCommunity(ctx, request.CommunityID)
	if err != nil {
		return nil, fmt.Errorf("DeleteCommunity: %w", err)
	}

	for _, id := range community.FollowerIDs {
		err = svc.db.UserRepo.UserDeleteCommunity(ctx, id, request.CommunityID)
		if err != nil {
			svc.log.Errorf("UserDeleteCommunity error: %s", err)
			return nil, err
		}
	}

	for _, id := range community.PostIDs {
		err = svc.db.PostRepo.DeletePost(ctx, id)
		if err != nil {
			svc.log.Errorf("DeletePost error: %s", err)
			return nil, err
		}

		err = svc.db.LikeRepo.DeleteLike(ctx, id)
		if err != nil {
			svc.log.Errorf("DeleteLike error: %s", err)
			return nil, err
		}
	}

	return &dto.DeleteCommunityResponse{}, nil
}

func NewCommunityService(log *logrus.Entry, db *db.Repository) CommunityService {
	return &communityServiceImpl{log: log, db: db}
}
