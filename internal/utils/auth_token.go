package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type AuthTokenWrapper struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

const (
	ViperJWTTTLKey    = "service.jwt_ttl"
	ViperJWTSecretKey = "service.jwt_secret"
)

func GenerateAuthToken(atw *AuthTokenWrapper) (string, error) {
	if atw.ExpiresAt == 0 {
		t := time.Second * time.Duration(viper.GetInt64(ViperJWTTTLKey))
		atw.ExpiresAt = time.Now().Add(t).Unix()
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atw)
	authToken, err := jwtToken.SignedString([]byte(viper.GetString(ViperJWTSecretKey)))
	if err != nil {
		return "", err
	}

	return authToken, nil
}
