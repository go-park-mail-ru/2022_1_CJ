package auth_common

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const saltSize = 32

func GetSalt() ([]byte, error) {
	passwordSalt := make([]byte, saltSize)
	_, err := rand.Read(passwordSalt)
	if err != nil {
		return passwordSalt, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}
	return passwordSalt, nil
}

func GetHash512(password string, salt []byte) ([]byte, error) {
	var passwordHash []byte
	sha512Hasher := sha512.New()
	if _, err := sha512Hasher.Write(append([]byte(password), salt...)); err != nil {
		return passwordHash, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}
	passwordHash = sha512Hasher.Sum(nil)
	return passwordHash, nil
}
