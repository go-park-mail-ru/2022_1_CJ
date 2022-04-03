package dto

// Only used in responses! Does not need validation.
type Post struct {
	AuthorID string   `json:"author_id"`
	PostID   string   `json:"post_id"`
	Message  string   `json:"message"`
	Images   []string `json:"images,omitempty"`
}

type GetPostDataRequest struct {
	Message string   `json:"message"`
	Images  []string `json:"images,omitempty"`
}

type GetPostDataResponse struct {
	Post Post `json:"post"`
}

type GetPostEditDataRequest struct {
	Message string   `json:"message"`
	Images  []string `json:"images,omitempty"`
}

type GetPostDeleteDataRequest struct{}

type GetPostRequest struct{}
