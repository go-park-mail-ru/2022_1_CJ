package utils

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCreateCookie(t *testing.T) {
	cookieGen := new(http.Cookie)
	cookieGen.Name = "Name"
	cookieGen.Value = "Value"
	cookieGen.Path = "/"
	cookieGen.Expires = time.Now().Add(time.Second * time.Duration(123))
	cookieGen.HttpOnly = false
	cookieGen.SameSite = http.SameSiteLaxMode
	cookieFunc := utils.CreateCookie("Name", "Value", 123)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, cookieFunc.Name, cookieGen.Name) {
			t.Error("got : ", cookieFunc.Name, " expected :", cookieGen.Name)
		}
		if !assert.Equal(t, cookieFunc.Value, cookieGen.Value) {
			t.Error("got : ", cookieFunc.Value, " expected :", cookieGen.Value)
		}
		if !assert.Equal(t, cookieFunc.HttpOnly, cookieGen.HttpOnly) {
			t.Error("got : ", cookieFunc.HttpOnly, " expected :", cookieGen.HttpOnly)
		}
	})
}

func TestCreateHttpOnlyCookie(t *testing.T) {
	cookieGen := new(http.Cookie)
	cookieGen.Name = "Name"
	cookieGen.Value = "Value"
	cookieGen.Path = "/"
	cookieGen.Expires = time.Now().Add(time.Second * time.Duration(123))
	cookieGen.HttpOnly = true
	cookieGen.SameSite = http.SameSiteLaxMode
	cookieFunc := utils.CreateHTTPOnlyCookie("Name", "Value", 123)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, cookieFunc.Name, cookieGen.Name) {
			t.Error("got : ", cookieFunc.Name, " expected :", cookieGen.Name)
		}
		if !assert.Equal(t, cookieFunc.Value, cookieGen.Value) {
			t.Error("got : ", cookieFunc.Value, " expected :", cookieGen.Value)
		}
		if !assert.Equal(t, cookieFunc.HttpOnly, cookieGen.HttpOnly) {
			t.Error("got : ", cookieFunc.HttpOnly, " expected :", cookieGen.HttpOnly)
		}
	})
}
