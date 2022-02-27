package constants

import (
	"errors"
	"net/http"
)

// COdedError is an error wrapper which wraps errors with http status codes.
type CodedError struct {
	err  error
	code int
}

func (ce *CodedError) Error() string {
	return ce.err.Error()
}

func (ce *CodedError) Code() int {
	return ce.code
}

var (
	// Unathorized
	ErrPasswordMismatch = &CodedError{errors.New("password mismatch"), http.StatusUnauthorized}

	// Bad Request
	ErrBindRequest     = &CodedError{errors.New("failed to bind request"), http.StatusBadRequest}
	ErrValidateRequest = &CodedError{errors.New("failed to validate request"), http.StatusBadRequest}
	ErrDBNotFound      = &CodedError{errors.New("not found in the database"), http.StatusBadRequest}

	// Internal
	ErrGenerateUUID = &CodedError{errors.New("failed to generate UUID"), http.StatusInternalServerError}

	// Conflict
	ErrEmailAlreadyTaken = &CodedError{errors.New("email is taken already by other user"), http.StatusConflict}
)
