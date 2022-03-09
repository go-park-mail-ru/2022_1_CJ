package dto

type Post struct {
	AuthorID string   `json:"author_id"`
	Message  string   `json:"message"`
	Images   []string `json:"images"`
}
