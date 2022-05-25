package controller

import (
	"context"

	auth_constants "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/constants"
	auth_db "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/handler"
	auth_dto "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/dto"
	auth_service "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/service"
	auth_utils "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServerImpl struct {
	log *logrus.Entry
	rep *auth_service.Registry
	handler.UnimplementedUserAuthServer
}

func CreateAuthServer(dbConn *mongo.Database, log *logrus.Entry) *AuthServerImpl {
	repository, err := auth_db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	registry := auth_service.NewRegistry(log, repository)
	return &AuthServerImpl{rep: registry, log: log}
}

func (s *AuthServerImpl) Login(ctx context.Context, in *handler.LoginReq) (*handler.LoginRes, error) {
	request := new(auth_dto.LoginUserRequest)

	request.Email = in.Email
	request.Password = in.Pwd

	response, err := s.rep.AuthService.LoginUser(context.Background(), request)
	if err != nil {
		s.log.Errorf("SignupUser error: %s", err)
		return &handler.LoginRes{}, err
	}

	return &handler.LoginRes{Token: response.AuthToken, UserID: response.UserID}, nil
}

func (s *AuthServerImpl) SignUp(ctx context.Context, in *handler.SignUpReq) (*handler.SignUpRes, error) {
	request := new(auth_dto.SignupUserRequest)

	request.Email = in.Email
	request.Password = in.Pwd

	response, err := s.rep.AuthService.SignupUser(context.Background(), request)
	if err != nil {
		s.log.Errorf("SignupUser: %s", err)
		return &handler.SignUpRes{}, err
	}

	return &handler.SignUpRes{Token: response.AuthToken, UserID: response.UserID}, nil
}

func (s *AuthServerImpl) Check(ctx context.Context, in *handler.CheckReq) (*handler.CheckRes, error) {
	res := &handler.CheckRes{
		Code:     false,
		NewToken: auth_constants.Nothing,
		UserID:   auth_constants.Nothing,
	}
	tw, err := auth_utils.ParseAuthToken(in.Token)
	if err != nil {
		return res, err
	}

	res.UserID = tw.UserID

	authToken, err := auth_utils.RefreshIfNeededAuthToken(tw)
	res.NewToken = authToken
	if err != nil {
		return res, err
	}

	if len(authToken) != 0 {
		res.Code = true
	}
	return res, nil
}
