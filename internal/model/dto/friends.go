package dto

type SendFriendRequestRequest struct {
	UserID string `json:"user_id"`
}

type SendFriendRequestResponse BasicResponse

type AcceptFriendRequestRequest struct {
	UserID     string `json:"user_id"`
	IsAccepted bool   `json:"is_accepted"`
}

type AcceptFriendRequestResponse struct {
	RequestsID []string `json:"request_ids"`
}

type DeleteFriendRequest struct {
	ExFriendID string `query:"ex_friend_id"`
}

type DeleteFriendResponse struct {
	FriendsID []string `json:"friend_ids"`
}

type GetFriendsRequests struct{}

type GetFriendsResponse struct {
	FriendsID []string `json:"friend_ids"`
}

type GetOutcomingRequestsRequest struct{}

type GetOutcomingRequestsResponse struct {
	RequestIDs []string `json:"request_ids"`
}

type GetIncomingRequestsRequest struct{}

type GetIncomingRequestsResponse struct {
	RequestIDs []string `json:"request_ids"`
}
