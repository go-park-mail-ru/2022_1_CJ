package core

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
)

// User describes a user entity
type User struct {
	ID           string          `bson:"_id"`
	Name         common.UserName `bson:"name"`
	Image        string          `bson:"images"`
	Email        string          `bson:"email"`
	Phone        string          `bson:"phone"`
	Location     string          `bson:"location"`
	BirthDay     string          `bson:"birth_day"`
	CreatedAt    int64           `bson:"created_at"` // unix timestamp
	Posts        []string        `bson:"posts,omitempty"`
	FriendsID    string          `bson:"friends_id"`
	DialogIDs    []string        `bson:"dialog_ids,omitempty"`
	CommunityIDs []string        `bson:"community_ids,omitempty"`
}

type EditInfo struct {
	Name     common.UserName `bson:"name"`
	Avatar   string          `bson:"avatar"`
	Phone    string          `bson:"phone"`
	Location string          `bson:"location"`
	BirthDay string          `bson:"birth_day"`
}
