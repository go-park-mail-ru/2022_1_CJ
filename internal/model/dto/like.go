package dto

// Only used in responses! Does not need validation.
type Like struct {
	Amount  int64    `json:"amount"`
	MyLike  bool     `json:"my_like"`
	UserIDs []string `json:"user_ids,omitempty"`
}

type IncreaseLikeRequest struct {
	PostID  string `json:"post_id,omitempty"`
	PhotoID string `json:"photo_id,omitempty"`
}

type IncreaseLikeResponse BasicResponse

type ReduceLikeRequest struct {
	PostID  string `json:"post_id,omitempty"`
	PhotoID string `json:"photo_id,omitempty"`
}

type ReduceLikeResponse BasicResponse

type GetLikePostRequest struct {
	PostID string `query:"post_id" validate:"required"`
}

type GetLikePostResponse struct {
	Likes Like `json:"likes"`
}
type GetLikePhotoRequest struct {
	PhotoID string `query:"photo_id" validate:"required"`
}

type GetLikePhotoResponse struct {
	Likes Like `json:"likes"`
}
