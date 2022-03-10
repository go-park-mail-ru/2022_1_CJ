package db

import (
	"errors"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"go.mongodb.org/mongo-driver/mongo"
)

// wrapError translates mongo's ErrNoDocuments to custom error.
func wrapError(err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return constants.ErrDBNotFound
	}
	return err
}
