package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

// Only used in responses! Does not need validation.
type User struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  common.UserName `json:"name"`
}

type GetUserDataRequest struct {
	UserID string `query:"user_id"`
}

type GetUserDataResponse struct {
	User User `json:"user"`
}

type GetUserPostsRequest struct {
	UserID string `query:"user_id"`
}

type GetUserPostsResponse struct {
	PostIDs []string `json:"post_ids"`
}

type GetUserFeedRequest struct{}

type GetUserFeedResponse struct {
	PostIDs []string `json:"post_ids"`
}
