package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
}

type GetUserDataRequest struct {
	UserID string `json:"user_id"`
}

type GetUserDataResponse struct {
	User User `json:"user"`
}

type GetUserFeedRequest struct {
	UserID string `json:"user_id"`
}

type GetUserFeedResponse struct {
	Posts []Post `json:"posts"`
}
