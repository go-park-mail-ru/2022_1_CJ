package dto

type Comment struct {
	ID      string   `json:"id"`
	Author  User     `json:"author"`
	Message string   `json:"message"`
	Images  []string `json:"images,omitempty"`
}

type CreateCommentRequest struct {
	PostID  string   `query:"post_id" validate:"required"`
	Message string   `json:"message" validate:"required"`
	Images  []string `json:"images,omitempty"`
}

type CreateCommentResponse BasicResponse

type GetCommentsRequest struct {
	PostID string `query:"post_id" validate:"required"`
	Limit  int64  `query:"limit,omitempty"`
	Page   int64  `query:"page,omitempty"`
}

type GetCommentsResponse struct {
	Comments    []Comment `json:"comments"`
	Total       int64     `json:"total"`
	AmountPages int64     `json:"amount_pages"`
}

type EditCommentRequest struct {
	PostID    string   `query:"post_id" validate:"required"`
	CommentID string   `query:"comment_id" validate:"required"`
	Message   string   `json:"message"`
	Images    []string `json:"images,omitempty"`
}

type EditCommentResponse BasicResponse

type DeleteCommentRequest struct {
	PostID    string `query:"post_id" validate:"required"`
	CommentID string `query:"comment_id" validate:"required"`
}

type DeleteCommentResponse struct{}
