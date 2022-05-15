// This file is in core bc GenUUID generates
// UUID that are needed for core models
package auth_core

import (
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/constants"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GenUUID() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", status.Error(codes.Internal, fmt.Errorf("%w: %v", auth_constants.ErrGenerateUUID, err).Error())
	}
	return uuid.String(), nil
}
