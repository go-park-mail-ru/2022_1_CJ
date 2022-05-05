package auth_utils

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type AuthTokenWrapper struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateAuthToken(atw *AuthTokenWrapper) (string, error) {
	if atw.ExpiresAt == 0 {
		t := time.Second * time.Duration(viper.GetInt64(auth_constants.ViperJWTTTLKey))
		atw.ExpiresAt = time.Now().Add(t).Unix()
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atw)
	authToken, err := jwtToken.SignedString([]byte(viper.GetString(auth_constants.ViperJWTSecretKey)))
	if err != nil {
		return "", status.Error(codes.Internal, auth_constants.ErrSignToken.Error())

	}
	return authToken, nil
}

func ParseAuthToken(authToken string) (*AuthTokenWrapper, error) {
	t, err := jwt.ParseWithClaims(
		authToken,
		&AuthTokenWrapper{},
		keyFunc([]byte(viper.GetString(auth_constants.ViperJWTSecretKey))),
	)

	if ve, ok := err.(*jwt.ValidationError); ok {
		// check if Expiration error was set
		if ve.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
			return nil, status.Error(codes.Internal, auth_constants.ErrAuthTokenExpired.Error())
		} else {
			return nil, status.Error(codes.Internal, auth_constants.ErrAuthTokenInvalid.Error())
		}
	} else if err != nil {
		return nil, status.Error(codes.Internal, auth_constants.ErrParseAuthToken.Error())
	}

	atw, ok := t.Claims.(*AuthTokenWrapper)
	if !ok {
		return nil, status.Error(codes.Internal, auth_constants.ErrAuthTokenInvalid.Error())
	}

	return atw, nil
}

func RefreshIfNeededAuthToken(atw *AuthTokenWrapper) (string, error) {
	if atw.ExpiresAt > time.Now().Unix()+viper.GetInt64(auth_constants.ViperJWTTTLKey)/2 {
		t := time.Second * time.Duration(viper.GetInt64(auth_constants.ViperJWTTTLKey))
		atw.ExpiresAt = time.Now().Add(t).Unix()
	} else {
		return "", nil
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atw)
	authToken, err := jwtToken.SignedString([]byte(viper.GetString(auth_constants.ViperJWTSecretKey)))
	if err != nil {
		return "", status.Error(codes.Internal, auth_constants.ErrSignToken.Error())
	}

	return authToken, nil
}

// KeyFunc returns key function for validating a token
func keyFunc(key []byte) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Internal, auth_constants.ErrUnexpectedSigningMethod.Error())
		}
		return key, nil
	}
}
