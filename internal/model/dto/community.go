package dto

// Only used in responses! Does not need validation.
type Community struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Only used in responses! Does not need validation.
type CommunityProfile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Info      string `json:"info"`
	Followers int64  `json:"followers"`
	Admins    []User `json:"admins"`
}

type GetCommunityRequest struct {
	CommunityID string `query:"community_id"`
}

type GetCommunityResponse struct {
	Community CommunityProfile `json:"community,omitempty"`
}

type GetCommunityPostsRequest struct {
	CommunityID string `query:"community_id"`
	Limit       int64  `query:"limit,omitempty"`
	Page        int64  `query:"page,omitempty"`
}

type GetCommunityPostsResponse struct {
	Posts       []GetPosts `json:"posts"`
	Total       int64      `json:"total"`
	AmountPages int64      `json:"amount_pages"`
}

type GetUserCommunitiesRequest struct {
	UserID string `query:"user_id"`
	Limit  int64  `query:"limit,omitempty"`
	Page   int64  `query:"page,omitempty"`
}

type GetUserCommunitiesResponse struct {
	Communities []Community `json:"communities,omitempty"`
	Total       int64       `json:"total"`
	AmountPages int64       `json:"amount_pages"`
}

type GetUserManageCommunitiesRequest struct {
	UserID string `query:"user_id"`
	Limit  int64  `query:"limit,omitempty"`
	Page   int64  `query:"page,omitempty"`
}

type GetUserManageCommunitiesResponse struct {
	Communities []Community `json:"communities,omitempty"`
	Total       int64       `json:"total"`
	AmountPages int64       `json:"amount_pages"`
}

type GetCommunitiesRequest struct {
	Limit int64 `query:"limit,omitempty"`
	Page  int64 `query:"page,omitempty"`
}

type GetCommunitiesResponse struct {
	Communities []Community `json:"communities,omitempty"`
	Total       int64       `json:"total"`
	AmountPages int64       `json:"amount_pages"`
}

type SearchCommunitiesRequest struct {
	Selector string `query:"selector" validate:"required"`
	Limit    int64  `query:"limit,omitempty"`
	Page     int64  `query:"page,omitempty"`
}

type SearchCommunitiesResponse struct {
	Communities []Community `json:"communities,omitempty"`
	Total       int64       `json:"total"`
	AmountPages int64       `json:"amount_pages"`
}

type UpdatePhotoCommunityRequest struct {
	CommunityID string `query:"community_id" validate:"required"`
}

type UpdatePhotoCommunityResponse struct {
	URL string `json:"url"`
}

type JoinCommunityRequest struct {
	CommunityID string `query:"community_id" validate:"required"`
}

type JoinCommunityResponse BasicResponse

type LeaveCommunityRequest struct {
	CommunityID string `query:"community_id" validate:"required"`
}

type LeaveCommunityResponse BasicResponse

type GetFollowersRequest struct {
	CommunityID string `query:"community_id" validate:"required"`
	Limit       int64  `query:"limit,omitempty"`
	Page        int64  `query:"page,omitempty"`
}

type GetFollowersResponse struct {
	Amount      int64  `json:"amount"`
	Followers   []User `json:"followers"`
	Total       int64  `json:"total"`
	AmountPages int64  `json:"amount_pages"`
}

type GetMutualFriendsRequest struct {
	CommunityID string `query:"community_id" validate:"required"`
	Limit       int64  `query:"limit,omitempty"`
	Page        int64  `query:"page,omitempty"`
}

type GetMutualFriendsResponse struct {
	Amount      int64  `json:"amount"`
	Followers   []User `json:"followers,omitempty"`
	Total       int64  `json:"total"`
	AmountPages int64  `json:"amount_pages"`
}

type CreateCommunityRequest struct {
	Name   string   `json:"name" validate:"required"`
	Image  string   `json:"image"`
	Info   string   `json:"info" validate:"required"`
	Admins []string `json:"admins"`
}

type CreateCommunityResponse BasicResponse

type EditCommunityRequest struct {
	CommunityID string   `json:"community_id"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Info        string   `json:"info"`
	Admins      []string `json:"admins"`
}

type EditCommunityResponse BasicResponse

type DeleteCommunityRequest struct {
	CommunityID string `query:"community_id"`
}

type DeleteCommunityResponse BasicResponse

type CreatePostCommunityRequest struct {
	CommunityID string   `json:"community_id"`
	Message     string   `json:"message" validate:"required"`
	Files       []string `json:"files,omitempty"`
}

type CreatePostCommunityResponse BasicResponse

type EditPostCommunityRequest struct {
	CommunityID string   `json:"community_id"`
	PostID      string   `json:"post_id"`
	Message     string   `json:"message"`
	Files       []string `json:"files,omitempty"`
}

type EditPostCommunityResponse BasicResponse

type DeletePostCommunityRequest struct {
	CommunityID string `query:"community_id"`
	PostID      string `query:"post_id"`
}

type DeletePostCommunityResponse struct{}
