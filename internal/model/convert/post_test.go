package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPost2DTOByUser(t *testing.T) {
	postCore := &core.Post{ID: "123", Message: "body", Images: []string{"img1", "img2"}}
	authorCore := &core.User{ID: "1234", Name: common.UserName{First: "Oleg", Last: "Krytoi"}, Image: "img3"}
	postDTO := Post2DTOByUser(postCore, authorCore)
	expect := dto.Post{ID: postCore.ID, Author: dto.Author{ID: authorCore.ID, Name: authorCore.Name.Full(), Image: authorCore.Image, Type: "User"},
		Message: postCore.Message, Images: postCore.Images}
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, postDTO, expect) {
			t.Error("got : ", postDTO, " expected :", expect)
		}
	})
}

func TestPost2DTOByCommunity(t *testing.T) {
	postCore := &core.Post{ID: "123", Message: "body", Images: []string{"img1", "img2"}}
	communityCore := &core.Community{ID: "1234", Name: "bestName", Image: "img3"}
	postDTO := Post2DTOByCommunity(postCore, communityCore)
	expect := dto.Post{ID: postCore.ID, Author: dto.Author{ID: communityCore.ID, Name: communityCore.Name, Image: communityCore.Image, Type: "Community"},
		Message: postCore.Message, Images: postCore.Images}
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, postDTO, expect) {
			t.Error("got : ", postDTO, " expected :", expect)
		}
	})
}
