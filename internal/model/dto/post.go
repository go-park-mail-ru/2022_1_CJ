package dto

// Only used in responses! Does not need validation.
type Post struct {
	AuthorID string   `json:"author_id"`
	PostID   string   `json:"post_id"`
	Message  string   `json:"message"`
	Images   []string `json:"images"`
}

type GetPostDataRequest struct {
	UserID  string   `json:"user_id"`
	Message string   `json:"message"`
	Images  []string `json:"images"`
}

type GetPostDataResponse struct {
	Post Post `json:"post"`
}

type GetPostEditDataRequest struct {
	UserID  string   `json:"user_id"`
	ID      string   `json:"post_id"`
	Message string   `json:"message"`
	Images  []string `json:"images"`
}

type GetPostDeleteDataRequest struct {
	UserID string `json:"user_id"`
	ID     string `json:"post_id"`
}
