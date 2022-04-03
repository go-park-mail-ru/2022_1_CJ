package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

type SignupUserRequest struct {
	Name     common.UserName `json:"name"`
	Email    string          `json:"email"    validate:"required,email"`
	Password string          `json:"password" validate:"required"`
}

type SignupUserResponse AuthTokenResponse
type SignupUserResponse BasicResponse

type LoginUserRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse AuthTokenResponse
