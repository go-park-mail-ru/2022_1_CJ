package service

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type FriendsService interface {
	SendFriendRequest(ctx context.Context, request *dto.SendFriendRequestRequest, userID string) (*dto.SendFriendRequestResponse, error)
	AcceptFriendRequest(ctx context.Context, request *dto.AcceptFriendRequestRequest, userID string) (*dto.AcceptFriendRequestResponse, error)
	DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest, userID string) (*dto.DeleteFriendResponse, error)
	GetFriendsByUserID(ctx context.Context, userID string) (*dto.GetFriendsResponse, error)
	GetFriendRequests(ctx context.Context, userID string) (*dto.GetRequestsResponse, error)
}

type friendsServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *friendsServiceImpl) SendFriendRequest(ctx context.Context, request *dto.SendFriendRequestRequest, userID string) (*dto.SendFriendRequestResponse, error) {
	if err := svc.db.FriendsRepo.IsUniqRequest(ctx, request.UserID, userID); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.IsNotFriend(ctx, request.UserID, userID); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.MakeRequest(ctx, request.UserID, userID); err != nil {
		return nil, err
	}

	return &dto.SendFriendRequestResponse{}, nil
}

// Проверить на самого себя!
func (svc *friendsServiceImpl) AcceptFriendRequest(ctx context.Context, request *dto.AcceptFriendRequestRequest, userID string) (*dto.AcceptFriendRequestResponse, error) {
	if request.IsAccepted {
		if err := svc.db.FriendsRepo.MakeFriends(ctx, userID, request.UserID); err != nil {
			return nil, err
		}
	}

	if err := svc.db.FriendsRepo.DeleteRequest(ctx, userID, request.UserID); err != nil {
		svc.log.Errorf("DeleteRequest error: %s", err)
		return nil, err
	}

	requests, err := svc.db.FriendsRepo.GetRequestsByUserID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetRequestsByUserID error: %s", err)
		return nil, err
	}

	return &dto.AcceptFriendRequestResponse{RequestsID: requests}, nil
}

func (svc *friendsServiceImpl) DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest, userID string) (*dto.DeleteFriendResponse, error) {
	if err := svc.db.FriendsRepo.DeleteFriend(ctx, userID, request.ExFriendID); err != nil {
		svc.log.Errorf("DeleteFriend error: %s", err)
		return nil, err
	}

	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetRequestsByUserID error: %s", err)
		return nil, err
	}

	return &dto.DeleteFriendResponse{FriendsID: friends}, nil
}

func (svc *friendsServiceImpl) GetFriendsByUserID(ctx context.Context, userID string) (*dto.GetFriendsResponse, error) {
	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetFriendsByUserID error: %s", err)
		return nil, err
	}
	return &dto.GetFriendsResponse{FriendsID: friends}, nil
}

func (svc *friendsServiceImpl) GetFriendRequests(ctx context.Context, userID string) (*dto.GetRequestsResponse, error) {
	requests, err := svc.db.FriendsRepo.GetRequestsByUserID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetRequestsByUserID error: %s", err)
		return nil, err
	}
	return &dto.GetRequestsResponse{RequestIDs: requests}, nil
}

func NewFriendsService(log *logrus.Entry, db *db.Repository) FriendsService {
	return &friendsServiceImpl{log: log, db: db}
}
