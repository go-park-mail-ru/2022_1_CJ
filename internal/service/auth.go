//go:generate mockgen -source=user.go -destination=user_mock.go -package=service
package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
)

type AuthService interface {
	SignupUser(ctx context.Context, request *dto.SignupUserRequest) (*dto.SignupUserResponse, error)
	LoginUser(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
}

type AuthServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *AuthServiceImpl) SignupUser(ctx context.Context, request *dto.SignupUserRequest) (*dto.SignupUserResponse, error) {
	if exists, err := svc.db.UserRepo.CheckUserEmailExistence(ctx, request.Email); err != nil {
		svc.log.Errorf("CheckUserEmailExistence error: %s", err)
		return nil, err
	} else if exists {
		svc.log.Errorf("CheckUserEmailExistence error: %s", constants.ErrEmailAlreadyTaken)
		return nil, constants.ErrEmailAlreadyTaken
	}
	user := &core.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := user.Password.Init(request.Password); err != nil {
		svc.log.Errorf("Init password error: %s", err)
		return nil, err
	}

	if err := svc.db.UserRepo.CreateUser(ctx, user); err != nil {
		svc.log.Errorf("CreateUser error: %s", err)
		return nil, err
	}

	if err := svc.db.FriendsRepo.CreateFriends(ctx, user.ID); err != nil {
		svc.log.Errorf("CreateFriends error: %s", err)
		return nil, err
	}
	svc.log.Debug("Create user success")

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		svc.log.Errorf("GenerateAuthToken error: %s", err)
		return nil, err
	}
	svc.log.Debugf("Generate auth token success; Token: %s", authToken)

	return &dto.SignupUserResponse{AuthToken: authToken}, nil
}

func (svc *AuthServiceImpl) LoginUser(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	user, err := svc.db.UserRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		svc.log.Errorf("GetUserByEmail error: %s", err)
		return nil, err
	}

	if err := user.Password.Validate(request.Password); err != nil {
		svc.log.Errorf("Validate error: %s", err)
		return nil, err
	}

	svc.log.Debug("Login success")

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		svc.log.Errorf("GenerateAuthToken error: %s", err)
		return nil, err
	}
	svc.log.Debugf("Generate auth token success; Token: %s", authToken)

	return &dto.LoginUserResponse{AuthToken: authToken}, nil
}

func NewAuthService(log *logrus.Entry, db *db.Repository) AuthService {
	return &AuthServiceImpl{log: log, db: db}
}
