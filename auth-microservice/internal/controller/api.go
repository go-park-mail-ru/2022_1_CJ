package handler

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/auth-microservice/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/auth-microservice/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/auth-microservice/internal/service"
	"github.com/go-park-mail-ru/2022_1_CJ/auth-microservice/internal/utils"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServerImpl struct {
	log *logrus.Entry
	rep *service.Registry
	UnimplementedUserAuthServer
}

func CreateAuthServer(dbConn *mongo.Database, log *logrus.Entry) UserAuthServer {
	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	registry := service.NewRegistry(log, repository)
	return &AuthServerImpl{rep: registry, log: log}
}

func (s *AuthServerImpl) Login(ctx context.Context, in *LoginReq) (*LoginRes, error) {

	return &LoginRes{}, nil
}

func (s *AuthServerImpl) Check(ctx echo.Context, in *CheckReq) (*CheckRes, error) {
	res := &CheckRes{
		Code:     false,
		Message:  constants.ErrParseToken,
		NewToken: constants.Nothing,
		UserID:   constants.Nothing,
	}
	tw, err := utils.ParseAuthToken(in.Token)
	if err != nil {
		return res, err
	}

	res.UserID = tw.UserID

	authToken, err := utils.RefreshIfNeededAuthToken(tw)
	res.NewToken = authToken
	if err != nil {
		res.Message = constants.ErrRefreshToken
		return res, err
	}

	if len(authToken) != 0 {
		res.Code = true
	}
	return res, nil
}
