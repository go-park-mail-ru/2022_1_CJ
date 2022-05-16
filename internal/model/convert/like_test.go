package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsMyLike(t *testing.T) {
	userIDs := []string{"123", "123124"}
	userIDF := "123"
	userIDS := "321"

	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, isMyLike(userIDs, userIDF), true) {
			t.Error("got : ", isMyLike(userIDs, userIDF), " expected :", true)
		}
		if !assert.Equal(t, isMyLike(userIDs, userIDS), false) {
			t.Error("got : ", isMyLike(userIDs, userIDF), " expected :", false)
		}
	})
}

func TestLike2DTO(t *testing.T) {
	like := &core.Like{Amount: 2, UserIDs: []string{"12312", "21312312"}}
	userID := "123"
	res := Like2DTO(like, userID)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, dto.Like{Amount: like.Amount, UserIDs: like.UserIDs, MyLike: false}) {
			t.Error("got : ", res, " expected :", dto.Like{Amount: like.Amount, UserIDs: like.UserIDs, MyLike: false})
		}
	})
}
