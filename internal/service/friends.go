package service

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type FriendsService interface {
	SendFriendRequest(ctx context.Context, request *dto.SendFriendRequestRequest, UserID string) (*dto.SendFriendRequestResponse, error)
	AcceptFriendRequest(ctx context.Context, request *dto.AcceptFriendRequestRequest, UserID string) (*dto.AcceptFriendRequestResponse, error)
	DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest, UserID string) (*dto.DeleteFriendResponse, error)
	GetFriendsByUserID(ctx context.Context, UserID string) (*dto.GetFriendsResponse, error)
	GetFriendRequests(ctx context.Context, UserID string) (*dto.GetRequestsResponse, error)
}

type friendsServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *friendsServiceImpl) SendFriendRequest(ctx context.Context, request *dto.SendFriendRequestRequest, UserID string) (*dto.SendFriendRequestResponse, error) {
	if err := svc.db.FriendsRepo.IsUniqRequest(ctx, request.UserID, UserID); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.IsNotFriend(ctx, request.UserID, UserID); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.MakeRequest(ctx, request.UserID, UserID); err != nil {
		return nil, err
	}

	return &dto.SendFriendRequestResponse{}, nil
}

// Проверить на самого себя!
func (svc *friendsServiceImpl) AcceptFriendRequest(ctx context.Context, request *dto.AcceptFriendRequestRequest, UserID string) (*dto.AcceptFriendRequestResponse, error) {
	if request.IsAccepted {
		if err := svc.db.FriendsRepo.MakeFriends(ctx, UserID, request.UserID); err != nil {
			return nil, err
		}
	}

	svc.log.Debug("MakeFriends success")

	if err := svc.db.FriendsRepo.DeleteRequest(ctx, UserID, request.UserID); err != nil {
		svc.log.Errorf("DeleteRequest error: %s", err)
		return nil, err
	}

	svc.log.Debug("DeleteRequest success")

	requests, err := svc.db.FriendsRepo.GetRequestsByUserID(ctx, UserID)
	if err != nil {
		svc.log.Errorf("GetRequestsByUserID error: %s", err)
		return nil, err
	}

	return &dto.AcceptFriendRequestResponse{RequestsID: requests}, nil
}

func (svc *friendsServiceImpl) DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest, UserID string) (*dto.DeleteFriendResponse, error) {
	if err := svc.db.FriendsRepo.DeleteFriend(ctx, UserID, request.ExFriendID); err != nil {
		svc.log.Errorf("DeleteFriend error: %s", err)
		return nil, err
	}

	svc.log.Debug("DeleteFriend success")

	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, UserID)
	if err != nil {
		svc.log.Errorf("GetRequestsByUserID error: %s", err)
		return nil, err
	}
	return &dto.DeleteFriendResponse{FriendsID: friends}, nil
}

func (svc *friendsServiceImpl) GetFriendsByUserID(ctx context.Context, UserID string) (*dto.GetFriendsResponse, error) {
	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, UserID)
	if err != nil {
		svc.log.Errorf("GetFriendsByUserID error: %s", err)
		return nil, err
	}
	return &dto.GetFriendsResponse{FriendsID: friends}, nil
}

func (svc *friendsServiceImpl) GetFriendRequests(ctx context.Context, UserID string) (*dto.GetRequestsResponse, error) {
	requests, err := svc.db.FriendsRepo.GetRequestsByUserID(ctx, UserID)
	if err != nil {
		svc.log.Errorf("GetRequestsByUserID error: %s", err)
		return nil, err
	}
	return &dto.GetRequestsResponse{RequestIDs: requests}, nil
}

func NewFriendsService(log *logrus.Entry, db *db.Repository) FriendsService {
	return &friendsServiceImpl{log: log, db: db}
}
