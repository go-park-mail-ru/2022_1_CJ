package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
}

type GetUserDataRequest struct {
	UserID string `json:"user_id"` // not required
}

type GetUserDataResponse struct {
	User User `json:"user"`
}

type GetUserFeedRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type GetUserFeedResponse struct {
	Posts []Post `json:"posts"`
}

// -------------------REQUEST
type ReqSendRequest struct {
	UserID   string `json:"user_id"`
	PersonID string `json:"person_id" validate:"required"`
}

// --------------------ACCEPT
type AcceptRequest struct {
	UserID     string `json:"user_id"`
	PersonID   string `json:"person_id" validate:"required"`
	IsAccepted bool   `json:"is_accepted" validate:"required"`
}

type AcceptResponse struct {
	RequestsID []string `json:"friends_id"`
}

// ------------------DELETE
type DeleteFriendRequest struct {
	UserID     string `json:"user_id"`
	ExFriendID string `json:"ex_friend_id" validate:"required"`
}

type DeleteFriendResponse struct {
	FriendsID []string `json:"friends_id"`
}
