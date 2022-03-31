package dto

type ReqSendRequest struct{}

type ReqSendResponse struct{}

type AcceptRequest struct {
	IsAccepted bool `json:"is_accepted" validate:"required"`
}

type AcceptResponse struct {
	RequestsID []string `json:"friends_id"`
}

type DeleteFriendRequest struct{}

type DeleteFriendResponse struct {
	FriendsID []string `json:"friends_id"`
}
