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
	Community CommunityProfile `json:"community"`
}

type GetCommunityPostsRequest struct {
	CommunityID string `query:"community_id"`
}

type GetCommunityPostsResponse struct {
	Posts []GetPosts `json:"posts"`
}

type GetUserCommunitiesRequest struct {
	UserID string `query:"user_id"`
}

type GetUserCommunitiesResponse struct {
	Communities []Community `json:"communities"`
}

type GetUserManageCommunitiesRequest struct {
	UserID string `query:"user_id"`
}

type GetUserManageCommunitiesResponse struct {
	Communities []Community `json:"communities"`
}

type GetCommunitiesRequest struct{}

type GetCommunitiesResponse struct {
	Communities []Community `json:"communities"`
}

type SearchCommunitiesRequest struct {
	Selector string `query:"selector" validate:"required"`
}

type SearchCommunitiesResponse struct {
	Communities []Community `json:"communities"`
}

type JoinCommunityRequest struct {
	CommunityID string `query:"community_id"`
}

type JoinCommunityResponse BasicResponse

type LeaveCommunityRequest struct {
	CommunityID string `query:"community_id"`
}

type LeaveCommunityResponse BasicResponse

type GetFollowersRequest struct {
	CommunityID string `query:"community_id"`
}

type GetFollowersResponse struct {
	Amount    int64  `json:"amount"`
	Followers []User `json:"followers"`
}

type GetMutualFriendsRequest struct {
	CommunityID string `query:"community_id"`
}

type GetMutualFriendsResponse struct {
	Amount    int64  `json:"amount"`
	Followers []User `json:"followers"`
}

type CreateCommunityRequest struct {
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Info   string   `json:"info"`
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
	Images      []string `json:"images,omitempty"`
}

type CreatePostCommunityResponse BasicResponse

type EditPostCommunityRequest struct {
	CommunityID string   `json:"community_id"`
	PostID      string   `json:"post_id"`
	Message     string   `json:"message"`
	Images      []string `json:"images,omitempty"`
}

type EditPostCommunityResponse BasicResponse

type DeletePostCommunityRequest struct {
	CommunityID string `query:"community_id"`
	PostID      string `query:"post_id"`
}

type DeletePostCommunityResponse struct{}
