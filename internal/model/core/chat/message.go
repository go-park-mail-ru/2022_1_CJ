package chat

type Dialog struct {
	ID         string   `bson:"_id"`
	AuthorIDs  []string `bson:"author_ids,omitempty" `
	MessageIDs []string `bson:"messages_ids,omitempty"`
}

type Message struct {
	ID        string `bson:"_id"`
	Text      string `bson:"text"`
	AuthorID  string `bson:"author_id"`
	CreatedAt int64  `bson:"created_at"`
}
