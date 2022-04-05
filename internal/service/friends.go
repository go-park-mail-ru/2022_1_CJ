package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type FriendsService interface {
	SendRequest(ctx context.Context, request *dto.ReqSendRequest, UserID string) (*dto.ReqSendResponse, error)
	AcceptRequest(ctx context.Context, request *dto.AcceptRequest, UserID string) (*dto.AcceptResponse, error)
	DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest, UserID string) (*dto.DeleteFriendResponse, error)
	// Getter
	GetFriends(ctx context.Context, UserID string) (*dto.GetFriendsResponse, error)
	GetRequests(ctx context.Context, UserID string) (*dto.GetRequestsResponse, error)
}

type friendsServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *friendsServiceImpl) SendRequest(ctx context.Context, request *dto.ReqSendRequest, UserID string) (*dto.ReqSendResponse, error) {
	if err := svc.db.FriendsRepo.IsUniqRequest(ctx, request.PersonID, UserID); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.IsNotFriend(ctx, request.PersonID, UserID); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.MakeRequest(ctx, request.PersonID, UserID); err != nil {
		return nil, err
	}

	return &dto.ReqSendResponse{}, nil
}

// Проверить на самого себя!
func (svc *friendsServiceImpl) AcceptRequest(ctx context.Context, request *dto.AcceptRequest, UserID string) (*dto.AcceptResponse, error) {
	if request.IsAccepted {
		if err := svc.db.FriendsRepo.MakeFriends(ctx, UserID, request.PersonID); err != nil {
			return nil, err
		}
	}

	if err := svc.db.FriendsRepo.DeleteRequest(ctx, UserID, request.PersonID); err != nil {
		return nil, err
	}

	requests, err := svc.db.FriendsRepo.GetRequestsByUserID(ctx, UserID)
	if err != nil {
		return nil, err
	}
	// Get Requests
	return &dto.AcceptResponse{RequestsID: requests}, nil
}

func (svc *friendsServiceImpl) DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest, UserID string) (*dto.DeleteFriendResponse, error) {
	if err := svc.db.FriendsRepo.DeleteFriend(ctx, UserID, request.ExFriendID); err != nil {
		return nil, err
	}

	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, UserID)
	if err != nil {
		return nil, err
	}
	// Get Requests
	return &dto.DeleteFriendResponse{FriendsID: friends}, nil
}

func (svc *friendsServiceImpl) GetFriends(ctx context.Context, UserID string) (*dto.GetFriendsResponse, error) {
	friends, err := svc.db.FriendsRepo.GetFriendsByUserID(ctx, UserID)
	if err != nil {
		return nil, err
	}
	// Get Requests
	return &dto.GetFriendsResponse{FriendsID: friends}, nil
}

func (svc *friendsServiceImpl) GetRequests(ctx context.Context, UserID string) (*dto.GetRequestsResponse, error) {
	requests, err := svc.db.FriendsRepo.GetRequestsByUserID(ctx, UserID)
	if err != nil {
		return nil, err
	}
	// Get Requests
	return &dto.GetRequestsResponse{RequestsID: requests}, nil
}

func NewFriendsService(log *logrus.Entry, db *db.Repository) FriendsService {
	return &friendsServiceImpl{log: log, db: db}
}
