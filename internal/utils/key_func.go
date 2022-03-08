package utils

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/golang-jwt/jwt"
)

// KeyFunc returns key function for validating a token
func keyFunc(key []byte) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrUnexpectedSigningMethod
		}
		return key, nil
	}
}
