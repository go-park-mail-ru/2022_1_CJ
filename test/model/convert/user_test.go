package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser2DTO(t *testing.T) {
	userCore := &core.User{ID: "13", Email: "cool@email.ru", Name: common.UserName{First: "first", Last: "Last"}, Image: "img1"}
	res := convert.User2DTO(userCore)
	expect := dto.User{ID: userCore.ID, Email: userCore.Email, Name: userCore.Name, Image: userCore.Image}
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, expect) {
			t.Error("got : ", res, " expected :", expect)
		}
	})
}

func TestUser2author(t *testing.T) {
	userDTO := dto.User{ID: "13", Email: "cool@email.ru", Name: common.UserName{First: "first", Last: "Last"}, Image: "img1"}
	res := convert.User2author(userDTO)
	expect := dto.Author{ID: userDTO.ID, Name: userDTO.Name.Full(), Image: userDTO.Image, Type: "User"}
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, expect) {
			t.Error("got : ", res, " expected :", expect)
		}
	})
}
