package utils

import (
	"net/http"
	"time"
)

func CreateCookie(name, value string, ttl int64) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(time.Second * time.Duration(ttl))
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
	return cookie
}
