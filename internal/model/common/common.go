package common

import "fmt"

type UserName struct {
	First string `json:"first" bson:"first"`
	Last  string `json:"last" bson:"last"`
}

func (un *UserName) Full() string {
	return fmt.Sprintf("%s %s", un.First, un.Last)
}

type PageResponse struct {
	Total       int64
	AmountPages int64
}
