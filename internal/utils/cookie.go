package utils

import (
	"net/http"
	"time"
)

func CreateCookie(name, value string, ttl int64) *http.Cookie {
	return createCookie(name, value, ttl, false)
}

func CreateHTTPOnlyCookie(name, value string, ttl int64) *http.Cookie {
	return createCookie(name, value, ttl, true)
}

func createCookie(name, value string, ttl int64, httpOnly bool) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(time.Second * time.Duration(ttl))
	cookie.Path = "/"
	cookie.HttpOnly = httpOnly
	cookie.SameSite = http.SameSiteLaxMode
	return cookie
}
