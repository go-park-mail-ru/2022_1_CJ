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
}

type postServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
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
