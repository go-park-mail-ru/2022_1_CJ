package cl

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	handler "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/handler"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type AuthRepository interface {
	Login(email, pass string) (string, string, error)
	Check(token string) (string, string, bool, error)
	SignUp(email, pass string) (string, string, error)
}

type AuthRepositoryImpl struct {
	log    *logrus.Entry
	client handler.UserAuthClient
}

func NewAuthRepository(log *logrus.Entry, cl handler.UserAuthClient) AuthRepository {
	return &AuthRepositoryImpl{log: log, client: cl}
}

func (redisConnect *AuthRepositoryImpl) ParseError(err error) error {
	getErr, ok := status.FromError(err)
	if ok {
		if val, ok := constants.ParseError[getErr.Message()]; ok {
			return val
		}
	}
	return fmt.Errorf(getErr.Message())
}

func (redisConnect *AuthRepositoryImpl) Login(email, pass string) (string, string, error) {
	res, err := redisConnect.client.Login(context.Background(), &handler.LoginReq{
		Email: email,
		Pwd:   pass,
	})
	if err != nil {
		return "", "", redisConnect.ParseError(err)

	}
	return res.Token, res.UserID, err
}

func (redisConnect *AuthRepositoryImpl) SignUp(email, pass string) (string, string, error) {
	redisConnect.log.Info(email, pass)
	res, err := redisConnect.client.SignUp(context.Background(), &handler.SignUpReq{
		Email: email,
		Pwd:   pass,
	})
	if err != nil {
		return "", "", redisConnect.ParseError(err)

	}
	return res.Token, res.UserID, err
}

func (redisConnect *AuthRepositoryImpl) Check(token string) (string, string, bool, error) {
	res, err := redisConnect.client.Check(context.Background(), &handler.CheckReq{Token: token})
	if err != nil {
		return "", "", false, redisConnect.ParseError(err)

	}
	return res.NewToken, res.UserID, res.Code, err
}
