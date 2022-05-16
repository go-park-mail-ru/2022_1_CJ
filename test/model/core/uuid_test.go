package core

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenUUID(t *testing.T) {
	uuidF, _ := core.GenUUID()
	uuidS, _ := core.GenUUID()
	t.Run("Check equals", func(t *testing.T) {
		if !assert.NotEqual(t, uuidF, uuidS) {
			t.Error("got : ", uuidF, " expected :", uuidS)
		}
	})
}
