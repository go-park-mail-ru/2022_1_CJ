package common

type UserName struct {
	First string `json:"first" bson:"first"`
	Last  string `json:"last" bson:"last"`
}

func (un *UserName) Full() string {
	return un.First + un.Last
}

type MessageInfo struct {
	Text     string `json:"text" bson:"text"`
	AuthorID string `json:"author_id" bson:"author_id"`
	DialogID string `json:"dialog_id" bson:"dialog_id"`
}

type DialogInfo struct {
	DialogID string `json:"dialog_id" bson:"dialog_id"`
	Title    string `json:"title" bson:"title"`
}
