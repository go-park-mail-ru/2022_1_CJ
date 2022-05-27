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

func TestPostNull(t *testing.T) *core.Post {
	t.Helper()

	return &core.Post{}
}

func TestCommentNull(t *testing.T) *core.Comment {
	t.Helper()

	return &core.Comment{}
}

func TestFriendsNull(t *testing.T) *core.Friends {
	t.Helper()

	return &core.Friends{}
}

func TestDialogNull(t *testing.T) *core.Dialog {
	t.Helper()

	return &core.Dialog{}
}

func TestFriends(t *testing.T) *core.Friends {
	t.Helper()
	return &core.Friends{
		ID:                "12345678",
		OutcomingRequests: []string{"123", "234"},
		IncomingRequests:  []string{"2345", "2357"},
		Friends:           []string{"123567", "213123", "214335345"},
	}
}

func TestPost(t *testing.T) *core.Post {
	t.Helper()
	return &core.Post{
		ID:        "12345678",
		AuthorID:  "123456789",
		Message:   "Hi it's my first post",
		Files:     []string{"src/image.jpg"},
		CreatedAt: 1323123,
	}
}

func TestComment(t *testing.T) *core.Comment {
	t.Helper()
	return &core.Comment{
		ID:        "12345678",
		AuthorID:  "123456789",
		Message:   "Hi it's my first post",
		Images:    []string{"src/image.jpg"},
		CreatedAt: 1323123,
	}
}
func TestLike(t *testing.T) *core.Like {
	t.Helper()
	return &core.Like{
		ID:        "1",
		Amount:    2,
		Subject:   "1",
		UserIDs:   []string{"1", "2"},
		CreatedAt: 1323123,
	}
}

func TestNullLike(t *testing.T) *core.Like {
	t.Helper()
	return &core.Like{}
}

func TestNullCommunity(t *testing.T) *core.Community {
	t.Helper()
	return &core.Community{}
}

func TestCommunity(t *testing.T) *core.Community {
	t.Helper()
	return &core.Community{
		ID:          "12345678",
		Name:        "Community",
		Image:       "image",
		Info:        "info",
		FollowerIDs: []string{"1", "2"},
		AdminIDs:    []string{"1", "2"},
		PostIDs:     []string{"1", "2"},
		CreatedAt:   124565,
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

func TestMessage(t *testing.T) *core.Message {
	t.Helper()
	return &core.Message{
		ID:        "12345678",
		Body:      "hi message",
		AuthorID:  "12345671",
		IsRead:    []core.IsRead{{Participant: "12345672", IsRead: false}},
		CreatedAt: 124565,
	}
}
