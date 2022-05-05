package core

type Community struct {
	ID          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Image       string   `bson:"image"`
	Info        string   `bson:"info"`
	FollowerIDs []string `bson:"followers,omitempty"`
	AdminIDs    []string `bson:"admins,omitempty"`
	PostIDs     []string `bson:"posts,omitempty"`
	CreatedAt   int64    `bson:"created_at"` // unix timestamp
}
