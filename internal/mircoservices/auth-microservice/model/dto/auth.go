package auth_dto

type SignupUserRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type SignupUserResponse struct {
	AuthToken string
	UserID    string
}

type LoginUserRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginUserResponse struct {
	AuthToken string
	UserID    string
}

type BasicResponse struct{}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
