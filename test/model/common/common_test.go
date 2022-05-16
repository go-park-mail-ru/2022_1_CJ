package common

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFull(t *testing.T) {
	name := common.UserName{First: "Alex", Last: "Cool"}
	fullName := name.Full()
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, fullName, name.First+" "+name.Last) {
			t.Error("got : ", fullName, " expected :", name.First+" "+name.Last)
		}
	})
}
