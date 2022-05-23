package dto

type SendFriendRequestRequest struct {
	From string `header:"User-Id" validate:"required"`
	To   string `json:"to" validate:"required"`
}

type SendFriendRequestResponse BasicResponse

type RevokeFriendRequestRequest struct {
	From string `header:"User-Id" validate:"required"`
	To   string `json:"to" validate:"required"`
}

type RevokeFriendRequestResponse BasicResponse

type AcceptFriendRequestRequest struct {
	To   string `header:"User-Id" validate:"required"`
	From string `json:"from" validate:"required"`
}

type AcceptFriendRequestResponse BasicResponse

type GetFriendsRequest struct {
	UserID      string `header:"User-Id" validate:"required"`
	QueryUserID string `query:"user_id"`
}

type GetFriendsResponse struct {
	FriendIDs []string `json:"friend_ids"`
}

type DeleteFriendRequest struct {
	UserID   string `header:"User-Id" validate:"required"`
	FriendID string `query:"friend_id" validate:"required"`
}

type DeleteFriendResponse BasicResponse

type GetIncomingRequestsRequest struct {
	UserID string `header:"User-Id" validate:"required"`
}

type GetIncomingRequestsResponse struct {
	RequestIDs []string `json:"request_ids"`
}

type GetOutcomingRequestsRequest struct {
	UserID string `header:"User-Id" validate:"required"`
}

type GetOutcomingRequestsResponse struct {
	RequestIDs []string `json:"request_ids"`
}
