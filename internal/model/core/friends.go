package core

type Friends struct {
	ID     string `bson:"_id"`
	UserID string `bson:"user_id"`

	Requests []string `bson:"requests,omitempty"`
	Friends  []string `bson:"friends,omitempty"`
}
