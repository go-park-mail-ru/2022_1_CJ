package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	UserRepo UserRepository
}

func NewRepository(dbConn *mongo.Database) (*Repository, error) {
	var err error
	repository := new(Repository)

	repository.UserRepo, err = NewUserRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	return repository, nil
}
