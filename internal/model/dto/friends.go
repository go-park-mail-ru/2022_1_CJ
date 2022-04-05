package dto

type ReqSendRequest struct {
	PersonID string `json:"person_id"`
}

type ReqSendResponse struct{}

type AcceptRequest struct {
	PersonID   string `json:"person_id"`
	IsAccepted bool   `json:"is_accepted" binding:"required"`
}

type AcceptResponse struct {
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
	RequestsID []string `json:"requests_id"`
}
