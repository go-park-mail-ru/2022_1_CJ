package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
}

// Add status
type UserProfile struct {
	UserInfo  User     `json:"user_info"`
	Avatar    string   `json:"avatar"`
	Phone     string   `json:"phone"`
	Location  string   `json:"location"`
	BirthDay  string   `json:"birth_day"`
	FriendIDs []string `json:"friend_ids"`
	PostIDs   []string `json:"post_ids"`
}

type EditProfile struct {
	Name     common.UserName `json:"name"`
	Avatar   string          `json:"avatar"`
	Phone    string          `json:"phone"`
	Location string          `json:"location"`
	BirthDay string          `json:"birth_day"`
}

type GetUserRequest struct {
	UserID string `query:"user_id"`
}

type GetUserResponse struct {
	User User `json:"user"`
}

type GetUserPostsRequest struct {
	UserID string `query:"user_id"`
}

type GetUserPostsResponse struct {
	PostIDs []string `json:"post_ids"`
}

type GetUserFeedResponse struct {
	PostIDs []string `json:"post_ids"`
}

type GetProfileRequest struct {
	UserID string `json:"user_id"`
}

type GetProfileResponse struct {
	UserProfile UserProfile `json:"user_profile"`
}

type EditProfileRequest struct {
	NewInfo EditProfile `json:"new_info"`
}

type EditProfileResponse struct {
	UserProfile UserProfile `json:"user_profile"`
}
