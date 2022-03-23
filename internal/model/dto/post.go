package dto

// Only used in responses! Does not need validation.
type Post struct {
	AuthorID string   `json:"author_id"`
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
