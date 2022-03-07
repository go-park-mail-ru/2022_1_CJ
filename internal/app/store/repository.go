package store

import "github.com/go-park-mail-ru/2022_1_CJ/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}
