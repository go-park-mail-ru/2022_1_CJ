package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type PostService interface {
	CreatePost(ctx context.Context, request *dto.GetPostDataRequest) (*dto.GetPostDataResponse, error)
	EditPost(ctx context.Context, request *dto.GetPostEditDataRequest) (*dto.GetPostDataResponse, error)
	DeletePost(ctx context.Context, request *dto.GetPostDeleteDataRequest) error
}

type postServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *postServiceImpl) DeletePost(ctx context.Context, request *dto.GetPostDeleteDataRequest) error {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return err
	}

	post, err := svc.db.PostRepo.GetPostByID(ctx, request.ID)
	if err != nil {
		return err
	}

	err = svc.db.PostRepo.DeletePost(ctx, post)
	if err != nil {
		return err
	}

	err = svc.db.UserRepo.UserDeletePost(ctx, user, request.ID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *postServiceImpl) EditPost(ctx context.Context, request *dto.GetPostEditDataRequest) (*dto.GetPostDataResponse, error) {
	_, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	post, err := svc.db.PostRepo.EditPost(ctx, &core.Post{
		AuthorID: request.ID,
		ID:       request.ID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		return nil, err
	}

	return &dto.GetPostDataResponse{Post: convert.Post2DTO(post)}, nil
}

func (svc *postServiceImpl) CreatePost(ctx context.Context, request *dto.GetPostDataRequest) (*dto.GetPostDataResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	post, err := svc.db.PostRepo.CreatePost(ctx, &core.Post{
		AuthorID: request.UserID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		return nil, err
	}

	err = svc.db.UserRepo.UserAddPost(ctx, user, post.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GetPostDataResponse{Post: convert.Post2DTO(post)}, nil
}

func NewPostService(log *logrus.Entry, db *db.Repository) PostService {
	return &postServiceImpl{log: log, db: db}
}
