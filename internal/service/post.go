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
	CreatePost(ctx context.Context, request *dto.GetPostDataRequest, UserID string) (*dto.GetPostDataResponse, error)
	EditPost(ctx context.Context, request *dto.GetPostEditDataRequest, UserID string, PostID string) (*dto.GetPostDataResponse, error)
	DeletePost(ctx context.Context, UserID string, PostID string) error
	GetPost(ctx context.Context, PostID string) (*dto.GetPostDataResponse, error)
}

type postServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *postServiceImpl) DeletePost(ctx context.Context, UserID string, PostID string) error {
	post, err := svc.db.PostRepo.GetPostByID(ctx, PostID)
	if err != nil {
		return err
	}

	err = svc.db.PostRepo.DeletePost(ctx, post)
	if err != nil {
		return err
	}

	err = svc.db.UserRepo.UserDeletePost(ctx, UserID, PostID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *postServiceImpl) EditPost(ctx context.Context, request *dto.GetPostEditDataRequest, UserID string, PostID string) (*dto.GetPostDataResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, UserID)
	if err != nil {
		return nil, err
	}

	err = svc.db.UserRepo.UserCheckPost(ctx, user, PostID)
	if err != nil {
		return nil, err
	}

	_, err = svc.db.PostRepo.GetPostByID(ctx, PostID)
	if err != nil {
		return nil, err
	}

	post, err := svc.db.PostRepo.EditPost(ctx, &core.Post{
		AuthorID: UserID,
		ID:       PostID,
		Message:  request.Message,
		Images:   request.Images,
	})

	if err != nil {
		return nil, err
	}

	return &dto.GetPostDataResponse{Post: convert.Post2DTO(post)}, nil
}

func (svc *postServiceImpl) CreatePost(ctx context.Context, request *dto.GetPostDataRequest, UserID string) (*dto.GetPostDataResponse, error) {
	post, err := svc.db.PostRepo.CreatePost(ctx, &core.Post{
		AuthorID: UserID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		return nil, err
	}

	err = svc.db.UserRepo.UserAddPost(ctx, UserID, post.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GetPostDataResponse{Post: convert.Post2DTO(post)}, nil
}

func (svc *postServiceImpl) GetPost(ctx context.Context, PostID string) (*dto.GetPostDataResponse, error) {

	post, err := svc.db.PostRepo.GetPostByID(ctx, PostID)
	if err != nil {
		return nil, err
	}

	return &dto.GetPostDataResponse{Post: convert.Post2DTO(post)}, nil
}

func NewPostService(log *logrus.Entry, db *db.Repository) PostService {
	return &postServiceImpl{log: log, db: db}
}
