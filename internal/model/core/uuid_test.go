package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenUUID(t *testing.T) {
	uuidF, _ := GenUUID()
	uuidS, _ := GenUUID()
	t.Run("Check equals", func(t *testing.T) {
		if !assert.NotEqual(t, uuidF, uuidS) {
			t.Error("got : ", uuidF, " expected :", uuidS)
		}
	})
}
