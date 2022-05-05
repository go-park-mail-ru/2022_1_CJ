package core

type Friends struct {
	ID                string   `bson:"_id"` // a user's ID
	OutcomingRequests []string `bson:"requests,omitempty"`
	IncomingRequest   []string `bson:"incoming_requests,omitempty"`
	Friends           []string `bson:"friends,omitempty"`
}
