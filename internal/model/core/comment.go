package core

type Comment struct {
	ID        string   `bson:"_id"`
	AuthorID  string   `bson:"author_id"`
	Message   string   `bson:"message"`
	Images    []string `bson:"images,omitempty"`
	CreatedAt int64    `bson:"created_at"` // unix timestamp
}
