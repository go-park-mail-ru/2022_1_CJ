package core

type IsRead struct {
	Participant string `bson:"_id" json:"id"`
	IsRead      bool   `bson:"is_read" json:"is_read"`
}

type Message struct {
	ID          string   `bson:"_id"`
	Body        string   `bson:"body"`
	AuthorID    string   `bson:"author_id"`
	IsRead      []IsRead `bson:"is_participants_read,omitempty"`
	Attachments []string `json:"attachments"`
	CreatedAt   int64    `bson:"created_at"` // unix timestamp
}

type Dialog struct {
	ID           string    `bson:"_id"`
	Name         string    `bson:"name"`
	Participants []string  `bson:"participants"`
	Messages     []Message `bson:"messages,omitempty"`
	CreatedAt    int64     `bson:"created_at"`
}
