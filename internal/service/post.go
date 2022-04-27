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
	GetPost(ctx context.Context, request *dto.GetPostRequest, userID string) (*dto.GetPostResponse, error)
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
		svc.log.Errorf("CreatePost error: %s", err)
		return nil, err
	}
	svc.log.Debug("CreatePost success")

	err = svc.db.UserRepo.UserAddPost(ctx, userID, post.ID)
	if err != nil {
		svc.log.Errorf("UserAddPost error: %s", err)
		return nil, err
	}

	_, err = svc.db.LikeRepo.CreateLike(ctx, &core.Like{Subject: post.ID})
	if err != nil {
		svc.log.Errorf("CreateLike error: %s", err)
		return nil, err
	}

	svc.log.Debugf("UserAddPost success; Current post ID: %s", post.ID)
	return &dto.CreatePostResponse{}, nil
}

func (svc *postServiceImpl) GetPost(ctx context.Context, request *dto.GetPostRequest, userID string) (*dto.GetPostResponse, error) {
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetPostByID error: %s", err)
		return nil, err
	}
	svc.log.Debug("GetPostByID success")

	author, err := svc.db.UserRepo.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetLikeBySubjectID error: %s", err)
		return nil, err
	}

	return &dto.GetPostResponse{Post: convert.Post2DTO(post, author), Likes: convert.Like2DTO(like, userID)}, nil
}

func (svc *postServiceImpl) EditPost(ctx context.Context, request *dto.EditPostRequest, userID string) (*dto.EditPostResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	err = svc.db.UserRepo.UserCheckPost(ctx, user, request.PostID)
	if err != nil {
		svc.log.Errorf("UserCheckPost error: %s", err)
		return nil, err
	}

	postBefore, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetPostByID error: %s", err)
		return nil, err
	}
	svc.log.Debugf("Post data befor edit: Message: %s; Images paths: %v", postBefore.Message, postBefore.Images)

	post, err := svc.db.PostRepo.EditPost(ctx, &core.Post{
		AuthorID: userID,
		ID:       request.PostID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		svc.log.Errorf("EditPost error: %s", err)
		return nil, err
	}

	svc.log.Debugf("Post data after edit: Message: %s; Images paths: %v", post.Message, post.Images)

	return &dto.EditPostResponse{}, nil
}

func (svc *postServiceImpl) DeletePost(ctx context.Context, request *dto.DeletePostRequest, userID string) (*dto.DeletePostResponse, error) {
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetPostByID error: %s", err)
		return nil, err
	}

	if post.AuthorID != userID {
		svc.log.Errorf("Not author error: %s", constants.ErrAuthorIDMismatch)
		return nil, constants.ErrAuthorIDMismatch
	}

	err = svc.db.PostRepo.DeletePost(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("DeletePost error: %s", err)
		return nil, err
	}
	svc.log.Debug("DeletePost success")
	err = svc.db.UserRepo.UserDeletePost(ctx, userID, request.PostID)
	if err != nil {
		svc.log.Errorf("UserDeletePost error: %s", err)
		return nil, err
	}
	svc.log.Debug("UserDeletePost success")
	return &dto.DeletePostResponse{}, nil
}

func NewPostService(log *logrus.Entry, db *db.Repository) PostService {
	return &postServiceImpl{log: log, db: db}
}
