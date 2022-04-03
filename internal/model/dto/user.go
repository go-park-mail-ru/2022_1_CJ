package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
}

type GetUserDataRequest struct{}

type GetUserDataResponse struct {
	User User `json:"user"`
}

type GetUserFeedRequest struct{}

type GetUserFeedResponse struct {
	PostsID []string `json:"post_ids"`
}
