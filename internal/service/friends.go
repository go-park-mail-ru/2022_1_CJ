package service

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type FriendsService interface {
	SendRequest(ctx context.Context, request *dto.SendFriendRequestRequest) (*dto.SendFriendRequestResponse, error)
	RevokeRequest(ctx context.Context, request *dto.RevokeFriendRequestRequest) (*dto.RevokeFriendRequestResponse, error)
	AcceptRequest(ctx context.Context, request *dto.AcceptFriendRequestRequest) (*dto.AcceptFriendRequestResponse, error)

	GetFriends(ctx context.Context, request *dto.GetFriendsRequest) (*dto.GetFriendsResponse, error)
	DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest) (*dto.DeleteFriendResponse, error)

	GetIncomingRequests(ctx context.Context, request *dto.GetIncomingRequestsRequest) (*dto.GetIncomingRequestsResponse, error)
	GetOutcomingRequests(ctx context.Context, request *dto.GetOutcomingRequestsRequest) (*dto.GetOutcomingRequestsResponse, error)
}

type friendsServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *friendsServiceImpl) SendRequest(ctx context.Context, request *dto.SendFriendRequestRequest) (*dto.SendFriendRequestResponse, error) {
	if err := svc.db.FriendsRepo.IsUniqRequest(ctx, request.From, request.To); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.IsNotFriend(ctx, request.From, request.To); err != nil {
		return nil, err
	}

	if err := svc.db.FriendsRepo.CreateRequest(ctx, request.From, request.To); err != nil {
		return nil, err
	}

	return &dto.SendFriendRequestResponse{}, nil
}

func (svc *friendsServiceImpl) RevokeRequest(ctx context.Context, request *dto.RevokeFriendRequestRequest) (*dto.RevokeFriendRequestResponse, error) {
	if err := svc.db.FriendsRepo.DeleteRequest(ctx, request.From, request.To); err != nil {
		return nil, err
	}
	return &dto.RevokeFriendRequestResponse{}, nil
}

func (svc *friendsServiceImpl) AcceptRequest(ctx context.Context, request *dto.AcceptFriendRequestRequest) (*dto.AcceptFriendRequestResponse, error) {
	if err := svc.db.FriendsRepo.DeleteRequest(ctx, request.From, request.To); err != nil {
		return nil, err
	}
	if err := svc.db.FriendsRepo.MakeFriends(ctx, request.From, request.To); err != nil {
		return nil, err
	}
	return &dto.AcceptFriendRequestResponse{}, nil
}

func (svc *friendsServiceImpl) GetFriends(ctx context.Context, request *dto.GetFriendsRequest) (*dto.GetFriendsResponse, error) {
	friendIDs, err := svc.db.FriendsRepo.GetFriends(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetFriendsResponse{FriendIDs: friendIDs}, nil
}

func (svc *friendsServiceImpl) DeleteFriend(ctx context.Context, request *dto.DeleteFriendRequest) (*dto.DeleteFriendResponse, error) {
	if err := svc.db.FriendsRepo.DeleteFriend(ctx, request.UserID, request.FriendID); err != nil {
		return nil, err
	}
	if err := svc.db.FriendsRepo.CreateRequest(ctx, request.FriendID, request.UserID); err != nil {
		return nil, err
	}
	return &dto.DeleteFriendResponse{}, nil
}

func (svc *friendsServiceImpl) GetIncomingRequests(ctx context.Context, request *dto.GetIncomingRequestsRequest) (*dto.GetIncomingRequestsResponse, error) {
	requests, err := svc.db.FriendsRepo.GetIncomingRequests(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetIncomingRequestsResponse{RequestIDs: requests}, nil
}

func (svc *friendsServiceImpl) GetOutcomingRequests(ctx context.Context, request *dto.GetOutcomingRequestsRequest) (*dto.GetOutcomingRequestsResponse, error) {
	requests, err := svc.db.FriendsRepo.GetOutcomingRequests(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.GetOutcomingRequestsResponse{RequestIDs: requests}, nil
}

func NewFriendsService(log *logrus.Entry, db *db.Repository) FriendsService {
	return &friendsServiceImpl{log: log, db: db}
}
