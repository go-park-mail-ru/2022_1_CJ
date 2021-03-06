package dto

type Author struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Type  string `json:"type"`
}

type Post struct {
	ID            string   `json:"id"`
	Author        Author   `json:"author"`
	Message       string   `json:"message"`
	Images        []string `json:"images,omitempty"`
	Attachments   []string `json:"attachments,omitempty"`
	CountComments int64    `json:"count_comments"`
	CreatedAt     int64    `json:"created_at"`
}

type CreatePostRequest struct {
	Message     string   `json:"message"`
	Images      []string `json:"images,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type CreatePostResponse BasicResponse

type GetPostRequest struct {
	PostID string `query:"post_id" validate:"required"`
}

type GetPostResponse struct {
	Post  Post `json:"post"`
	Likes Like `json:"likes"`
}

type EditPostRequest struct {
	PostID      string   `json:"post_id"`
	Message     string   `json:"message"`
	Images      []string `json:"images,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type EditPostResponse BasicResponse

type DeletePostRequest struct {
	PostID string `query:"post_id"`
}

type DeletePostResponse struct{}
