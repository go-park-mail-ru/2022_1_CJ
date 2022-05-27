//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
)

type OAuthService interface {
	AuthenticateThroughTelergam(ctx context.Context, request *dto.AuthenticateThroughTelergamRequest) (*dto.AuthenticateThroughTelergamResponse, error)
}

type OAuthServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *OAuthServiceImpl) AuthenticateThroughTelergam(ctx context.Context, request *dto.AuthenticateThroughTelergamRequest) (*dto.AuthenticateThroughTelergamResponse, error) {
	// TODO: check hash

	exists, err := svc.db.UserRepo.CheckUserIDExistence(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("CheckUserIDExistence: %w", err)
	}

	csrfToken, err := utils.GenerateCSRFToken(request.ID)
	if err != nil {
		return nil, fmt.Errorf("GenerateCSRFToken: %w", err)
	}

	if !exists {
		user := &core.User{
			ID:    request.ID,
			Name:  common.UserName{First: request.FirstName, Last: request.LastName},
			Image: request.PhotoURL,
		}

		if err := svc.db.UserRepo.InsertUser(ctx, user); err != nil {
			return nil, fmt.Errorf("CreateUser: %w", err)
		}

		if err := svc.db.FriendsRepo.CreateFriends(ctx, user.ID); err != nil {
			return nil, fmt.Errorf("CreateFriends: %w", err)
		}
	}

	return &dto.AuthenticateThroughTelergamResponse{CSRFToken: csrfToken}, nil
}

func NewOAuthService(log *logrus.Entry, db *db.Repository) OAuthService {
	return &OAuthServiceImpl{log: log, db: db}
}
