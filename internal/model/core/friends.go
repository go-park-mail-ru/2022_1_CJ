package core

type Friends struct {
	ID       string   `bson:"_id"` // a user's ID
	Requests []string `bson:"requests,omitempty"`
	Friends  []string `bson:"friends,omitempty"`
}
