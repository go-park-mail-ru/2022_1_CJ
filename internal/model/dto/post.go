package dto

type Post struct {
	ID      string   `json:"id"`
	Author  User     `json:"author"`
	Message string   `json:"message"`
	Images  []string `json:"images,omitempty"`
}

type CreatePostRequest struct {
	Message string   `json:"message" validate:"required"`
	Images  []string `json:"images,omitempty"`
}

type CreatePostResponse BasicResponse

type GetPostRequest struct {
	PostID string `query:"post_id" validate:"required"`
}

type GetPostResponse struct {
	Post Post `json:"post"`
}

type EditPostRequest struct {
	PostID  string   `json:"post_id"`
	Message string   `json:"message"`
	Images  []string `json:"images,omitempty"`
}

type EditPostResponse BasicResponse

type DeletePostRequest struct {
	PostID string `query:"post_id"`
}

type DeletePostResponse struct{}
