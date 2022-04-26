package db

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"testing"
)

func TestUser(t *testing.T) *core.User {
	t.Helper()

	return &core.User{
		ID: "123456789",
		Name: common.UserName{
			First: "Sasha",
			Last:  "Userov",
		},
		Email: "user@example.org",
		Image: "src/img.jpg",
		Phone: "+8(800)-555-35-35",
	}
}

func TestUserNull(t *testing.T) *core.User {
	t.Helper()

	return &core.User{}
}

func TestPost(t *testing.T) *core.Post {
	t.Helper()
	return &core.Post{
		ID:        "12345678",
		AuthorID:  "123456789",
		Message:   "Hi it's my first post",
		Images:    []string{"src/image.jpg"},
		CreatedAt: 1323123,
	}
}

func TestDialog(t *testing.T) *core.Dialog {
	t.Helper()
	return &core.Dialog{
		ID:           "12345678",
		Name:         "My dialog",
		Participants: []string{"12345671", "12345672"},
		Messages:     []core.Message{},
		CreatedAt:    124565,
	}
}
