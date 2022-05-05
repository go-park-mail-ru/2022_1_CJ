package auth_service

import (
	authdb "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/db"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	AuthService AuthService
}

func NewRegistry(log *logrus.Entry, repository *authdb.Repository) *Registry {
	registry := new(Registry)

	registry.AuthService = NewAuthService(log, repository)
	return registry
}
