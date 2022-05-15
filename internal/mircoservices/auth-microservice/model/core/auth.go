package auth_core

import (
	"encoding/base64"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserPassword describes user's password data
type UserPassword struct {
	Hash string `bson:"hash"`
	Salt string `bson:"salt"`
}

// User describes a user entity
type User struct {
	ID        string       `bson:"_id"`
	Email     string       `bson:"email"`
	Password  UserPassword `bson:"password"`
	CreatedAt int64        `bson:"created_at"` // unix timestamp
}

// Init generates salt and hash with given password and fills corresponding fields.
func (up *UserPassword) Init(password string) error {
	salt, err := auth_common.GetSalt()
	if err != nil {
		return status.Error(codes.Internal, auth_constants.ErrPasswordSalt.Error())
	}

	hash, err := auth_common.GetHash512(password, salt)
	if err != nil {
		return status.Error(codes.Internal, auth_constants.ErrPassword.Error())
	}

	up.Salt = base64.URLEncoding.EncodeToString(salt)
	up.Hash = base64.URLEncoding.EncodeToString(hash)

	return nil
}

// Validate checks if the given password is the one that is stored.
func (up *UserPassword) Validate(password string) error {
	salt, err := base64.URLEncoding.DecodeString(up.Salt)
	if err != nil {
		return status.Error(codes.Internal, fmt.Errorf("error decoding user's salt: %s", err).Error())
	}

	hash, err := auth_common.GetHash512(password, salt)
	if err != nil {
		return status.Error(codes.Internal, fmt.Errorf("error generating hash: %s", err).Error())
	}

	if base64.URLEncoding.EncodeToString(hash) != up.Hash {
		return status.Error(codes.Internal, auth_constants.ErrPasswordMismatch.Error())
	}

	return nil
}
