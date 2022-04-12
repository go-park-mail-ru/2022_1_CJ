package db

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"testing"
)

func TestUser(t *testing.T) *core.User {
	t.Helper()

	return &core.User{
		Name: common.UserName{
			First: "Sasha",
			Last:  "Userov",
		},
		Email: "user@example.org",
		Image: "src/img.jpg",
		Phone: "+8(800)-555-35-35",
	}
}
