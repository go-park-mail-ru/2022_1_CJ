package dto

type ReqSendRequest struct{}

type ReqSendResponse struct{}

type AcceptRequest struct {
	IsAccepted bool `json:"is_accepted" binding:"required"`
}

type AcceptResponse struct {
	RequestsID []string `json:"requests_id"`
}

type DeleteFriendRequest struct{}

type DeleteFriendResponse struct {
	FriendsID []string `json:"friends_id"`
}

type GetFriendsRequests struct{}

type GetFriendsResponse struct {
	FriendsID []string `json:"friends_id"`
}

type GetRequestsRequests struct{}

type GetRequestsResponse struct {
	RequestsID []string `json:"requests_id"`
}
