package core

type Post struct {
	AuthorID string   `bson:"author_id"`
	Message  string   `bson:"message"`
	Images   []string `bson:"images"`
}
