package core

type Like struct {
	ID        string   `bson:"_id"`
	Subject   string   `bson:"subject_id"`
	Amount    int64    `bson:"amount"`
	UserIDs   []string `bson:"user_ids,omitempty"`
	CreatedAt int64    `bson:"created_at"` // unix timestamp
}
