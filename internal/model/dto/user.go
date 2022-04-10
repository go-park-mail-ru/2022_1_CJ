package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
	Image string          `json:"image"`
}

// Add status
type UserProfile struct {
	ID        string          `json:"id"`
	Email     string          `json:"email"`
	Name      common.UserName `json:"name"`
	Avatar    string          `json:"avatar"`
	Phone     string          `json:"phone"`
	Location  string          `json:"location"`
	BirthDay  string          `json:"birth_day"`
	FriendIDs []string        `json:"friend_ids"`
	PostIDs   []string        `json:"post_ids"`
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
	Posts []Post `json:"posts"`
}

type GetUserFeedRequest struct{}

type GetUserFeedResponse struct {
	Posts []Post `json:"posts"`
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
