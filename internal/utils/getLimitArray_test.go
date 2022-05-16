package utils

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLimitArray(t *testing.T) {
	arrString := []string{"1", "2", "3", "4", "5"}
	res, total, page := GetLimitArray(&arrString, 2, 1)
	t.Run("Check pagination", func(t *testing.T) {
		if !assert.Equal(t, res, []string{"5", "4"}) {
			t.Error("got : ", res, " expected :", []string{"1", "2"})
		}
		if !assert.Equal(t, int(total), len(arrString)) {
			t.Error("got : ", total, " expected :", len(arrString))
		}
		if !assert.Equal(t, int(page), 3) {
			t.Error("got : ", total, " expected :", 3)
		}
	})
}

func TestGetLimitMessage(t *testing.T) {
	arrMessage := []core.Message{{ID: "123"}, {ID: "234"}, {ID: "23123"}, {ID: "12312"}, {ID: "231"}}
	res, total, page := GetLimitMessage(&arrMessage, 2, 1)
	t.Run("Check pagination", func(t *testing.T) {
		if !assert.Equal(t, res, []core.Message{{ID: "231"}, {ID: "12312"}}) {
			t.Error("got : ", res, " expected :", []core.Message{{ID: "231"}, {ID: "12312"}})
		}
		if !assert.Equal(t, int(total), len(arrMessage)) {
			t.Error("got : ", total, " expected :", len(arrMessage))
		}
		if !assert.Equal(t, int(page), 3) {
			t.Error("got : ", total, " expected :", 3)
		}
	})
}

func TestIsLarge(t *testing.T) {
	res := IsLarge(true)
	t.Run("Check res", func(t *testing.T) {
		if !assert.Equal(t, 1, int(res)) {
			t.Error("got : ", res, " expected :", 1)
		}
	})
}

func TestReverseString(t *testing.T) {
	arrString := []string{"1", "2", "3"}
	res := reverseString(arrString)
	t.Run("Check res", func(t *testing.T) {
		if !assert.Equal(t, res, []string{"3", "2", "1"}) {
			t.Error("got : ", res, " expected :", []string{"3", "2", "1"})
		}
	})
}

func TestReverseMessage(t *testing.T) {
	arrMessage := []core.Message{{ID: "123"}, {ID: "234"}, {ID: "23123"}}
	res := reverseMessage(arrMessage)
	t.Run("Check res", func(t *testing.T) {
		if !assert.Equal(t, res, []core.Message{{ID: "23123"}, {ID: "234"}, {ID: "123"}}) {
			t.Error("got : ", res, " expected :", []core.Message{{ID: "23123"}, {ID: "234"}, {ID: "123"}})
		}
	})
}
