package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComment2DTO(t *testing.T) {
	commentCore := &core.Comment{ID: "123", Message: "body", Images: []string{"img1", "img2"}}
	authorCore := &core.User{ID: "1234", Name: common.UserName{First: "Oleg", Last: "Krytoi"}, Image: "img3"}

	result := Comment2DTO(commentCore, authorCore)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, result, dto.Comment{ID: "123", Message: "body", Images: []string{"img1", "img2"}, Author: User2DTO(authorCore)}) {
			t.Error("got : ", result, " expected :", dto.Comment{ID: "123", Message: "body", Images: []string{"img1", "img2"}, Author: User2DTO(authorCore)})
		}
	})
}
