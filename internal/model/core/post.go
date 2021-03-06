package core

type Post struct {
	ID          string   `bson:"_id"`
	AuthorID    string   `bson:"author_id"`
	Message     string   `bson:"message"`
	Images      []string `bson:"images,omitempty"`
	Attachments []string `bson:"attachments,omitempty"`
	CreatedAt   int64    `bson:"created_at"` // unix timestamp
	Type        string   `bson:"type"`
	CommentsIDs []string `bson:"comment_ids,omitempty"`
}
