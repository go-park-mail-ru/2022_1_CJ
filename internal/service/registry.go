package service

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	AuthService AuthService
}

func NewRegistry(log *logrus.Entry, repository *db.Repository) *Registry {
	registry := new(Registry)

	registry.AuthService = NewAuthService(log, repository)

	return registry
}