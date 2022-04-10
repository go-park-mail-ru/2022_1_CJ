package common

type Message struct {
	ID         string `bson:"_id" json:"_id"`
	Text       string `bson:"text" json:"text"`
	ToUserID   string `bson:"to_user_id" json:"to_user_id"`
	FromUserID string `bson:"from_user_id" json:"from_user_id"`
	CreatedAt  int64  `bson:"created_at" json:"created_at"`
}
