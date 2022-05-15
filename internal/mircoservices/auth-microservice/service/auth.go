package auth_service

import (
	"context"
	authconstants "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/constants"
	authdb "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/db"
	authcore "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/core"
	authdto "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/dto"
	authutils "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	SignupUser(ctx context.Context, request *authdto.SignupUserRequest) (*authdto.SignupUserResponse, error)
	LoginUser(ctx context.Context, request *authdto.LoginUserRequest) (*authdto.LoginUserResponse, error)
}

type AuthServiceImpl struct {
	log *logrus.Entry
	db  *authdb.Repository
}

func (svc *AuthServiceImpl) LoginUser(ctx context.Context, request *authdto.LoginUserRequest) (*authdto.LoginUserResponse, error) {
	user, err := svc.db.AuthRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		svc.log.Errorf("GetUserByEmail error: %s", err)
		return nil, err
	}

	if err := user.Password.Validate(request.Password); err != nil {
		svc.log.Errorf("Validate error: %s", err)
		return nil, err
	}

	// AUTH
	authToken, err := authutils.GenerateAuthToken(&authutils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		svc.log.Errorf("GenerateAuthToken error: %s", err)
		return nil, err
	}

	return &authdto.LoginUserResponse{AuthToken: authToken, UserID: user.ID}, nil
}

func (svc *AuthServiceImpl) SignupUser(ctx context.Context, request *authdto.SignupUserRequest) (*authdto.SignupUserResponse, error) {
	if exists, err := svc.db.AuthRepo.CheckUserEmailExistence(ctx, request.Email); err != nil {
		return nil, err
	} else if exists {
		svc.log.Errorf("CheckUserEmailExistence error: %s", authconstants.ErrEmailAlreadyTaken)
		return nil, status.Error(codes.Internal, authconstants.ErrEmailAlreadyTaken.Error())
	}
	user := &authcore.User{
		Email: request.Email,
	}

	if err := user.Password.Init(request.Password); err != nil {
		svc.log.Errorf("Init password error: %s", err)
		return nil, err
	}

	id, err := svc.db.AuthRepo.CreateUser(ctx, user)
	if err != nil {
		svc.log.Errorf("CreateUser error: %s", err)
		return nil, err
	}
	// AUTH
	authToken, err := authutils.GenerateAuthToken(&authutils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		svc.log.Errorf("GenerateAuthToken error: %s", err)
		return nil, err
	}
	return &authdto.SignupUserResponse{AuthToken: authToken, UserID: id}, nil
}

func NewAuthService(log *logrus.Entry, db *authdb.Repository) AuthService {
	return &AuthServiceImpl{log: log, db: db}
}
