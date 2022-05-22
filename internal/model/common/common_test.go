package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFull(t *testing.T) {
	name := UserName{First: "Alex", Last: "Cool"}
	fullName := name.Full()
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, fullName, name.First+" "+name.Last) {
			t.Error("got : ", fullName, " expected :", name.First+" "+name.Last)
		}
	})
}
