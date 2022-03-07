package core

import (
	"encoding/base64"
	"fmt"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
)

// UserPassword describes user's password data
type UserPassword struct {
	Hash string `bson:"hash"`
	Salt string `bson:"salt"`
}

// User describes a user entity
type User struct {
	ID   string          `bson:"_id"`
	Name common.UserName `bson:"name"`

	Email string `bson:"email"`
	Phone string `bson:"phone"`

	CreatedAt int64 `bson:"created_at"`           // unix timestamp
	UpdatedAt int64 `bson:"updated_at,omitempty"` // unix timestamp

	Password UserPassword `bson:"password"`
}

// Init generates salt and hash with given password and fills corresponding fields.
func (up *UserPassword) Init(password string) error {
	salt, err := common.GetSalt()
	if err != nil {
		return fmt.Errorf("error generating salt: %s", err)
	}

	hash, err := common.GetHash512(password, salt)
	if err != nil {
		return fmt.Errorf("error generating hash: %s", err)
	}

	up.Salt = base64.URLEncoding.EncodeToString(salt)
	up.Hash = base64.URLEncoding.EncodeToString(hash)

	return nil
}

// Validate checks if the given password is the one that is stored.
func (up *UserPassword) Validate(password string) error {
	salt, err := base64.URLEncoding.DecodeString(up.Salt)
	if err != nil {
		return fmt.Errorf("error decoding user's salt: %s", err)
	}

	hash, err := common.GetHash512(password, salt)
	if err != nil {
		return fmt.Errorf("error generating hash: %s", err)
	}

	if base64.URLEncoding.EncodeToString(hash) != up.Hash {
		return constants.ErrPasswordMismatch
	}

	return nil
}