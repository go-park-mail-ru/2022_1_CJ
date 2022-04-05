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

type PostService interface {
	CreatePost(ctx context.Context, request *dto.CreatePostRequest, userID string) (*dto.CreatePostResponse, error)
	GetPost(ctx context.Context, request *dto.GetPostRequest) (*dto.GetPostResponse, error)
	EditPost(ctx context.Context, request *dto.EditPostRequest, userID string) (*dto.EditPostResponse, error)
	DeletePost(ctx context.Context, request *dto.DeletePostRequest, userID string) (*dto.DeletePostResponse, error)
}

type postServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *postServiceImpl) CreatePost(ctx context.Context, request *dto.CreatePostRequest, userID string) (*dto.CreatePostResponse, error) {
	post, err := svc.db.PostRepo.CreatePost(ctx, &core.Post{
		AuthorID: userID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		return nil, err
	}

	err = svc.db.UserRepo.UserAddPost(ctx, userID, post.ID)
	if err != nil {
		return nil, err
	}

	return &dto.CreatePostResponse{Post: convert.Post2DTO(post)}, nil
}

func (svc *postServiceImpl) GetPost(ctx context.Context, request *dto.GetPostRequest) (*dto.GetPostResponse, error) {
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		return nil, err
	}
	return &dto.GetPostResponse{Post: convert.Post2DTO(post)}, nil
}

func (svc *postServiceImpl) EditPost(ctx context.Context, request *dto.EditPostRequest, userID string) (*dto.EditPostResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = svc.db.UserRepo.UserCheckPost(ctx, user, request.PostID)
	if err != nil {
		return nil, err
	}

	_, err = svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		return nil, err
	}

	post, err := svc.db.PostRepo.EditPost(ctx, &core.Post{
		AuthorID: userID,
		ID:       request.PostID,
		Message:  request.Message,
		Images:   request.Images,
	})

	if err != nil {
		return nil, err
	}

	return &dto.EditPostResponse{Post: convert.Post2DTO(post)}, nil
}

func (svc *postServiceImpl) DeletePost(ctx context.Context, request *dto.DeletePostRequest, userID string) (*dto.DeletePostResponse, error) {
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != userID {
		return nil, constants.ErrAuthorIDMismatch
	}

	err = svc.db.PostRepo.DeletePost(ctx, request.PostID)
	if err != nil {
		return nil, err
	}

	err = svc.db.UserRepo.UserDeletePost(ctx, userID, request.PostID)
	if err != nil {
		return nil, err
	}

	return &dto.DeletePostResponse{}, nil
}

func NewPostService(log *logrus.Entry, db *db.Repository) PostService {
	return &postServiceImpl{log: log, db: db}
}
