package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

type RegisterUserRequest struct {
	Name     common.UserName `json:"name"`
	Email    string          `json:"email"    validate:"required,email"`
	Password string          `json:"password" validate:"required"`
}

type RegisterUserResponse BasicResponse
