package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfile2DTO(t *testing.T) {
	userCore := &core.User{ID: "13", Email: "cool@email.ru", Name: common.UserName{First: "first", Last: "Last"}}
	res := Profile2DTO(userCore)
	expect := dto.UserProfile{ID: userCore.ID, Email: userCore.Email, Name: userCore.Name}
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, expect) {
			t.Error("got : ", res, " expected :", expect)
		}
	})
}
