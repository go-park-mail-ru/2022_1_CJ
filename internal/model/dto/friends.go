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
	RequestsID []string `json:"requests_id"`
}

type DeleteFriendRequest struct {
	ExFriendID string `json:"ex_friend_id"`
}

type DeleteFriendResponse struct {
	FriendsID []string `json:"friends_id"`
}

type GetFriendsRequests struct{}

type GetFriendsResponse struct {
	FriendsID []string `json:"friends_id"`
}

type GetRequestsRequests struct{}

type GetRequestsResponse struct {
	RequestIDs []string `json:"request_ids"`
}
