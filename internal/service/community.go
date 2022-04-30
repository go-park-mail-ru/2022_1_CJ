package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type CommunityService interface {
	CreateCommunity(ctx context.Context, request *dto.CreateCommunityRequest, userID string) (*dto.CreateCommunityResponse, error)
	DeleteCommunity(ctx context.Context, request *dto.DeleteCommunityRequest, userID string) (*dto.DeleteCommunityResponse, error)
	EditCommunity(ctx context.Context, request *dto.EditCommunityRequest, userID string) (*dto.EditCommunityResponse, error)
	GetCommunity(ctx context.Context, request *dto.GetCommunityRequest) (*dto.GetCommunityResponse, error)
	GetCommunityPosts(ctx context.Context, request *dto.GetCommunityPostsRequest, userID string) (*dto.GetCommunityPostsResponse, error)
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
	err = svc.db.UserRepo.UserAddCommunity(ctx, userID, community.ID)
	if err != nil {
		svc.log.Errorf("UserAddCommunity error: %s", err)
		return nil, err
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

func (svc *communityServiceImpl) GetCommunityPosts(ctx context.Context, request *dto.GetCommunityPostsRequest, userID string) (*dto.GetCommunityPostsResponse, error) {
	community, err := svc.db.CommunityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		svc.log.Errorf("GetCommunityByID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	var posts []dto.GetPosts
	for _, id := range community.PostIDs {
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
		var admins []dto.User
		for _, id := range community.AdminIDs {
			user, err := svc.db.UserRepo.GetUserByID(ctx, id)
			if err != nil {
				svc.log.Errorf("GetCommunityByID error: %s", err)
				return nil, constants.ErrDBNotFound
			}
			admins = append(admins, convert.User2DTO(user))
		}
		posts = append(posts, dto.GetPosts{Post: convert.Post2DTOByCommunity(post, community, admins), Likes: convert.Like2DTO(like, userID)})
	}

	return &dto.GetCommunityPostsResponse{}, nil
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
		svc.log.Errorf("DeleteCommunity error: %s", err)
		return nil, err
	}

	for _, id := range community.FollowerIDs {
		err = svc.db.UserRepo.UserDeleteCommunity(ctx, id, request.CommunityID)
		if err != nil {
			svc.log.Errorf("UserDeleteCommunity error: %s", err)
			return nil, err
		}
	}

	return &dto.DeleteCommunityResponse{}, nil
}

func NewCommunityService(log *logrus.Entry, db *db.Repository) CommunityService {
	return &communityServiceImpl{log: log, db: db}
}
