package auth_constants

import (
	"errors"
	"net/http"
)

// CodedError is an error wrapper which wraps errors with http status codes.
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
	ErrMissingAuthToken  = &CodedError{errors.New("missing authorization token"), http.StatusUnauthorized}
	ErrMissingAuthCookie = &CodedError{errors.New("missing authorization cookie"), http.StatusUnauthorized}

	ErrPasswordMismatch = &CodedError{errors.New("password mismatch"), http.StatusUnauthorized}

	ErrAuthTokenInvalid        = &CodedError{errors.New("authorization token is invalid"), http.StatusUnauthorized}
	ErrUnexpectedSigningMethod = &CodedError{errors.New("unexpected signing method"), http.StatusUnauthorized}

	// Forbidden
	ErrAuthTokenExpired = &CodedError{errors.New("authorization token is expired"), http.StatusForbidden}
	ErrAuthorIDMismatch = &CodedError{errors.New("author id mismatch"), http.StatusForbidden}

	// Bad Request
	ErrValidateRequest = &CodedError{errors.New("failed to validate request"), http.StatusBadRequest}
	ErrDBNotFound      = &CodedError{errors.New("not found in the database"), http.StatusBadRequest}
	ErrPassword        = &CodedError{errors.New("error generating hash"), http.StatusBadRequest}
	ErrPasswordSalt    = &CodedError{errors.New("error generating salt"), http.StatusBadRequest}

	// Internal
	ErrSignToken      = &CodedError{errors.New("failed to sign token"), http.StatusInternalServerError}
	ErrGenerateUUID   = &CodedError{errors.New("failed to generate UUID"), http.StatusInternalServerError}
	ErrParseAuthToken = &CodedError{errors.New("failed to parse authorization token"), http.StatusInternalServerError}

	// Conflict
	ErrEmailAlreadyTaken = &CodedError{errors.New("email is taken already by other user"), http.StatusConflict}
)
