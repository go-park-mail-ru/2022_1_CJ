package core

type Friends struct {
	ID                string   `bson:"_id"` // User ID
	Friends           []string `bson:"friends,omitempty"`
	IncomingRequests  []string `bson:"incoming_requests,omitempty"`
	OutcomingRequests []string `bson:"outcoming_requests,omitempty"`
}
