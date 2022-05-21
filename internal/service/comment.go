package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/sirupsen/logrus"
)

type CommentService interface {
	CreateComment(ctx context.Context, request *dto.CreateCommentRequest, userID string) (*dto.CreateCommentResponse, error)
	GetComments(ctx context.Context, request *dto.GetCommentsRequest) (*dto.GetCommentsResponse, error)
	EditComment(ctx context.Context, request *dto.EditCommentRequest, userID string) (*dto.EditCommentResponse, error)
	DeleteComment(ctx context.Context, request *dto.DeleteCommentRequest, userID string) (*dto.DeleteCommentResponse, error)
}

type CommentServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *CommentServiceImpl) CreateComment(ctx context.Context, request *dto.CreateCommentRequest, userID string) (*dto.CreateCommentResponse, error) {
	Comment, err := svc.db.CommentRepo.CreateComment(ctx, &core.Comment{
		AuthorID: userID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		svc.log.Errorf("CreateComment error: %s", err)
		return nil, err
	}

	err = svc.db.PostRepo.PostAddComment(ctx, userID, Comment.ID)
	if err != nil {
		svc.log.Errorf("UserAddComment error: %s", err)
		return nil, err
	}

	return &dto.CreateCommentResponse{}, nil
}

func (svc *CommentServiceImpl) GetComments(ctx context.Context, request *dto.GetCommentsRequest) (*dto.GetCommentsResponse, error) {
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetPostByID error: %s", err)
		return nil, err
	}

	commentsIDs, total, pages := utils.GetLimitArray(&post.CommentsIDs, request.Limit, request.Page)
	var comments []dto.Comment
	for _, id := range commentsIDs {
		comment, err := svc.db.CommentRepo.GetCommentByID(ctx, id)
		if err != nil {
			svc.log.Errorf("GetCommunityByID error: %s", err)
			return nil, constants.ErrDBNotFound
		}
		user, err := svc.db.UserRepo.GetUserByID(ctx, comment.AuthorID)
		if err != nil {
			svc.log.Errorf("GetUserByID error: %s", err)
			return nil, constants.ErrDBNotFound
		}
		comments = append(comments, convert.Comment2DTO(comment, user))
	}

	return &dto.GetCommentsResponse{Comments: comments, AmountPages: pages, Total: total}, nil
}

func (svc *CommentServiceImpl) EditComment(ctx context.Context, request *dto.EditCommentRequest, userID string) (*dto.EditCommentResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	comment, err := svc.db.CommentRepo.GetCommentByID(ctx, request.CommentID)
	if err != nil {
		svc.log.Errorf("UserCheckComment error: %s", err)
		return nil, err
	}

	if user.ID != comment.AuthorID {
		svc.log.Errorf("Not author error: %s", constants.ErrAuthorIDMismatch)
		return nil, constants.ErrAuthorIDMismatch
	}

	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetPostByID error: %s", err)
		return nil, err
	}

	err = svc.db.PostRepo.PostCheckComment(ctx, post, request.CommentID)
	if err != nil {
		svc.log.Errorf("PostCheckComment error: %s", err)
		return nil, err
	}

	_, err = svc.db.CommentRepo.EditComment(ctx, &core.Comment{
		AuthorID: userID,
		ID:       request.CommentID,
		Message:  request.Message,
		Images:   request.Images,
	})
	if err != nil {
		svc.log.Errorf("EditComment error: %s", err)
		return nil, err
	}

	return &dto.EditCommentResponse{}, nil
}

func (svc *CommentServiceImpl) DeleteComment(ctx context.Context, request *dto.DeleteCommentRequest, userID string) (*dto.DeleteCommentResponse, error) {
	user, err := svc.db.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		svc.log.Errorf("GetUserByID error: %s", err)
		return nil, err
	}

	comment, err := svc.db.CommentRepo.GetCommentByID(ctx, request.CommentID)
	if err != nil {
		svc.log.Errorf("UserCheckComment error: %s", err)
		return nil, err
	}

	if user.ID != comment.AuthorID {
		svc.log.Errorf("Not author error: %s", constants.ErrAuthorIDMismatch)
		return nil, constants.ErrAuthorIDMismatch
	}

	post, err := svc.db.PostRepo.GetPostByID(ctx, request.PostID)
	if err != nil {
		svc.log.Errorf("GetPostByID error: %s", err)
		return nil, err
	}

	err = svc.db.PostRepo.PostCheckComment(ctx, post, request.CommentID)
	if err != nil {
		svc.log.Errorf("PostCheckComment error: %s", err)
		return nil, err
	}

	err = svc.db.CommentRepo.DeleteComment(ctx, request.CommentID)
	if err != nil {
		svc.log.Errorf("DeleteComment error: %s", err)
		return nil, err
	}

	err = svc.db.PostRepo.PostDeleteComment(ctx, request.PostID, request.CommentID)
	if err != nil {
		svc.log.Errorf("DeleteComment error: %s", err)
		return nil, err
	}

	return &dto.DeleteCommentResponse{}, nil
}

func NewCommentService(log *logrus.Entry, db *db.Repository) CommentService {
	return &CommentServiceImpl{log: log, db: db}
}
