package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type LikeService interface {
	IncreaseLike(ctx context.Context, request *dto.IncreaseLikeRequest, userID string) (*dto.IncreaseLikeResponse, error)
	ReduceLike(ctx context.Context, request *dto.ReduceLikeRequest, userID string) (*dto.ReduceLikeResponse, error)
	GetLikePost(ctx context.Context, request *dto.GetLikePostRequest, userID string) (*dto.GetLikePostResponse, error)
	GetLikePhoto(ctx context.Context, request *dto.GetLikePhotoRequest, userID string) (*dto.GetLikePhotoResponse, error)
}

type likeServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *likeServiceImpl) IncreaseLike(ctx context.Context, request *dto.IncreaseLikeRequest, userID string) (*dto.IncreaseLikeResponse, error) {
	if (request.PhotoID == "" && request.PostID == "") || (request.PhotoID != "" && request.PostID != "") {
		svc.log.Errorf("ErrBadJson error")
		return nil, constants.ErrBadJson
	}
	var subject string
	if request.PhotoID != "" {
		subject = request.PhotoID
	} else {
		subject = request.PostID
	}

	_, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, subject)
	if err != nil {
		svc.log.Errorf("GetLikeBySubjectID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	err = svc.db.LikeRepo.IncreaseLike(ctx, subject, userID)
	if err != nil {
		svc.log.Errorf("IncreaseLike error: %s", err)
		return nil, err
	}

	return &dto.IncreaseLikeResponse{}, nil
}
func (svc *likeServiceImpl) ReduceLike(ctx context.Context, request *dto.ReduceLikeRequest, userID string) (*dto.ReduceLikeResponse, error) {
	if (request.PhotoID == "" && request.PostID == "") || (request.PhotoID != "" && request.PostID != "") {
		svc.log.Errorf("ErrBadJson error")
		return nil, constants.ErrBadJson
	}
	var subject string
	if request.PhotoID != "" {
		subject = request.PhotoID
	} else {
		subject = request.PostID
	}

	_, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, subject)
	if err != nil {
		svc.log.Errorf("GetLikeBySubjectID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	err = svc.db.LikeRepo.ReduceLike(ctx, subject, userID)
	if err != nil {
		svc.log.Errorf("IncreaseLike error: %s", err)
		return nil, err
	}

	return &dto.ReduceLikeResponse{}, nil
}
func (svc *likeServiceImpl) GetLikePost(ctx context.Context, request *dto.GetLikePostRequest, userID string) (*dto.GetLikePostResponse, error) {
	like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetLikeBySubjectID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	return &dto.GetLikePostResponse{Likes: convert.Like2DTO(like, userID)}, nil
}

func (svc *likeServiceImpl) GetLikePhoto(ctx context.Context, request *dto.GetLikePhotoRequest, userID string) (*dto.GetLikePhotoResponse, error) {
	like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, request.PhotoID)
	if err != nil {
		svc.log.Errorf("GetLikeBySubjectID error: %s", err)
		return nil, constants.ErrDBNotFound
	}

	return &dto.GetLikePhotoResponse{Likes: convert.Like2DTO(like, userID)}, nil
}

func NewLikeService(log *logrus.Entry, db *db.Repository) LikeService {
	return &likeServiceImpl{log: log, db: db}
}
