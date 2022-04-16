package common

type UserName struct {
	First string `json:"first" bson:"first"`
	Last  string `json:"last" bson:"last"`
}

func (un *UserName) Full() string {
	return un.First + un.Last
}
