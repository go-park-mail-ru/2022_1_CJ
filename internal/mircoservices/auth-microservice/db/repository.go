package auth_db

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	AuthRepo AuthRepository
}

func NewRepository(dbConn *mongo.Database) (*Repository, error) {
	var err error
	repository := new(Repository)

	repository.AuthRepo, err = NewAuthRepository(dbConn)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("create repository %s", err).Error())
	}
	return repository, nil
}
