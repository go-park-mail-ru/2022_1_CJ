package service

import (
	"context"
	"fmt"

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
		AuthorID:    userID,
		Message:     request.Message,
		Images:      request.Images,
		Attachments: request.Attachments,
		Type:        constants.UserPost,
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
	postCore, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		return nil, err
	}

	var post dto.Post
	switch postCore.Type {
	case constants.UserPost:
		author, errUser := svc.db.UserRepo.GetUserByID(ctx, postCore.AuthorID)
		if errUser != nil {
			return nil, err
		}
		post = convert.Post2DTOByUser(postCore, author)

	case constants.CommunityPost:
		community, errComm := svc.db.CommunityRepo.GetCommunityByID(ctx, postCore.AuthorID)
		if errComm != nil {
			return nil, err
		}
		post = convert.Post2DTOByCommunity(postCore, community)
	default:
		return nil, constants.ErrDBNotFound
	}

	like, err := svc.db.LikeRepo.GetLikeBySubjectID(ctx, request.PostID)
	if err != nil {
		return nil, err
	}

	return &dto.GetPostResponse{Post: post, Likes: convert.Like2DTO(like, userID)}, nil
}

func (svc *postServiceImpl) EditPost(ctx context.Context, request *dto.EditPostRequest, userID string) (*dto.EditPostResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	err = svc.db.UserRepo.UserCheckPost(ctx, user, request.PostID)
	if err != nil {
		return nil, fmt.Errorf("UserCheckPost: %w", err)
	}

	postBefore, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		return nil, fmt.Errorf("GetPostByID: %w", err)
	}

	if len(request.Message) != 0 {
		postBefore.Message = request.Message
	}

	if request.Images != nil {
		postBefore.Images = request.Images
	}

	if request.Attachments != nil {
		postBefore.Attachments = request.Attachments
	}

	_, err = svc.db.PostRepo.EditPost(ctx, postBefore)
	if err != nil {
		return nil, fmt.Errorf("EditPost: %w", err)
	}

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

	err = svc.db.UserRepo.UserDeletePost(ctx, userID, request.PostID)
	if err != nil {
		svc.log.Errorf("UserDeletePost error: %s", err)
		return nil, err
	}

	err = svc.db.LikeRepo.DeleteLike(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("DeleteLike error: %s", err)
		return nil, err
	}

	return &dto.DeletePostResponse{}, nil
}

func NewPostService(log *logrus.Entry, db *db.Repository) PostService {
	return &postServiceImpl{log: log, db: db}
}
