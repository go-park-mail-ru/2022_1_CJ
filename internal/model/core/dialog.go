package core

type Message struct {
	ID        string `bson:"_id"`
	Body      string `bson:"body"`
	AuthorID  string `bson:"author_id"`
	CreatedAt int64  `bson:"created_at"` // unix timestamp
}

type Dialog struct {
	ID           string    `bson:"_id"`
	Name         string    `bson:"name"`
	Participants []string  `bson:"participants"`
	Messages     []Message `bson:"messages"`
}
