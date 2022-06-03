package dto

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
)

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
	Image string          `json:"image"`
}

// Add status
type UserProfile struct {
	ID       string          `json:"id"`
	Email    string          `json:"email"`
	Name     common.UserName `json:"name"`
	Avatar   string          `json:"avatar"`
	Phone    string          `json:"phone"`
	Location string          `json:"location"`
	BirthDay string          `json:"birth_day"`
}

type EditProfile struct {
	Name     common.UserName `json:"name"`
	Avatar   string          `json:"avatar"`
	Phone    string          `json:"phone"`
	Location string          `json:"location"`
	BirthDay string          `json:"birth_day"`
}

type GetPosts struct {
	Post  Post `json:"post"`
	Likes Like `json:"likes"`
}

type GetUserRequest struct {
	UserID string `query:"user_id"`
}

type GetUserResponse struct {
	User User `json:"user"`
}

type GetUserPostsRequest struct {
	UserID string `query:"user_id"`
	Limit  int64  `query:"limit,omitempty"`
	Page   int64  `query:"page,omitempty"`
}

type GetUserPostsResponse struct {
	Posts       []GetPosts `json:"posts"`
	Total       int64      `json:"total"`
	AmountPages int64      `json:"amount_pages"`
}

type GetUserFeedRequest struct {
	UserID               string `header:"User-Id" validate:"required"`
	PaginationParameters core.PaginationParameters
}

type GetUserFeedResponse struct {
	Posts       []GetPosts `json:"posts"`
	Total       int64      `json:"total"`
	AmountPages int64      `json:"amount_pages"`
}

type GetProfileRequest struct {
	UserID string `query:"user_id"`
}

type GetProfileResponse struct {
	UserProfile UserProfile `json:"user_profile"`
}

type EditProfileRequest struct {
	Name     common.UserName `json:"name"`
	Avatar   string          `json:"avatar"`
	Phone    string          `json:"phone"`
	Location string          `json:"location"`
	BirthDay string          `json:"birth_day"`
}

type EditProfileResponse BasicResponse

type UpdatePhotoRequest struct{}

type UpdatePhotoResponse struct {
	URL string `json:"url"`
}

type SearchUsersRequest struct {
	Selector string `query:"selector" validate:"required"`
	Limit    int64  `query:"limit,omitempty"`
	Page     int64  `query:"page,omitempty"`
}

type SearchUsersResponse struct {
	Users       []User `json:"users"`
	Total       int64  `json:"total"`
	AmountPages int64  `json:"amount_pages"`
}
