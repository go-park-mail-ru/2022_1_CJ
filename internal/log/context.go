package log

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/sirupsen/logrus"
)

func WithFieldsFromContext(ctx context.Context, log *logrus.Entry) *logrus.Entry {
	if value, ok := ctx.Value(constants.CtxKeyUserID{}).(string); ok {
		log = log.WithField("user_id", value)
	}

	if value, ok := ctx.Value(constants.CtxKeyXRequestID{}).(string); ok {
		log = log.WithField("x_request_id", value)
	}

	return log
}
