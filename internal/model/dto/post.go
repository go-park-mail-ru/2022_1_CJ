package dto

type Post struct {
	AuthorID string   `json:"author_id"`
	PostID   string   `json:"post_id"`
	Message  string   `json:"message"`
	Images   []string `json:"images,omitempty"`
}

type CreatePostRequest struct {
	Message string   `json:"message" validate:"required"`
	Images  []string `json:"images,omitempty"`
}

type CreatePostResponse struct {
	Post Post `json:"post"`
}

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

type EditPostResponse struct {
	Post Post `json:"post"`
}

type DeletePostRequest struct {
	PostID string `query:"post_id"`
}

type DeletePostResponse struct{}
