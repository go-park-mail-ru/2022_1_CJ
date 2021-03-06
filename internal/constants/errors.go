package constants

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

	ErrMissingCSRFCookie = &CodedError{errors.New("missing csrf cookie"), http.StatusUnauthorized}
	ErrCSRFTokenWrong    = &CodedError{errors.New("wrong csrf token in cookie"), http.StatusUnauthorized}

	ErrPasswordMismatch = &CodedError{errors.New("password mismatch"), http.StatusUnauthorized}

	ErrAuthTokenInvalid        = &CodedError{errors.New("authorization token is invalid"), http.StatusUnauthorized}
	ErrUnexpectedSigningMethod = &CodedError{errors.New("unexpected signing method"), http.StatusUnauthorized}

	ErrHashInvalid = &CodedError{errors.New("hash is invalid"), http.StatusUnauthorized}

	// Forbidden
	ErrAuthTokenExpired = &CodedError{errors.New("authorization token is expired"), http.StatusForbidden}
	ErrAuthorIDMismatch = &CodedError{errors.New("author id mismatch"), http.StatusForbidden}

	// Bad Request
	ErrBindRequest     = &CodedError{errors.New("failed to bind request"), http.StatusBadRequest}
	ErrValidateRequest = &CodedError{errors.New("failed to validate request"), http.StatusBadRequest}
	ErrDBNotFound      = &CodedError{errors.New("not found in the database"), http.StatusBadRequest}
	ErrBadJson         = &CodedError{errors.New("bad json request"), http.StatusBadRequest}
	ErrPassword        = &CodedError{errors.New("error generating hash"), http.StatusBadRequest}

	// Internal
	ErrSignToken      = &CodedError{errors.New("failed to sign token"), http.StatusInternalServerError}
	ErrGenerateUUID   = &CodedError{errors.New("failed to generate UUID"), http.StatusInternalServerError}
	ErrParseAuthToken = &CodedError{errors.New("failed to parse authorization token"), http.StatusInternalServerError}

	// Conflict
	ErrEmailAlreadyTaken = &CodedError{errors.New("email is taken already by other user"), http.StatusConflict}

	// Not Uniq
	ErrAddYourself         = &CodedError{errors.New("can't make yourself friend"), http.StatusConflict}
	ErrRequestAlreadyExist = &CodedError{errors.New("your request already was sent"), http.StatusConflict}
	ErrAlreadyFriends      = &CodedError{errors.New("your already friend with this person"), http.StatusConflict}
	ErrAlreadyFollower     = &CodedError{errors.New("you already in community"), http.StatusConflict}

	// Chat
	ErrSingleChat         = &CodedError{errors.New("you can't create dialog with no one"), http.StatusBadRequest}
	ErrDialogAlreadyExist = &CodedError{errors.New("dialog already exist"), http.StatusConflict}
)

var (
	ParseError = map[string]*CodedError{
		ErrMissingAuthToken.Error():        ErrMissingAuthToken,
		ErrMissingAuthCookie.Error():       ErrMissingAuthCookie,
		ErrMissingCSRFCookie.Error():       ErrMissingCSRFCookie,
		ErrCSRFTokenWrong.Error():          ErrCSRFTokenWrong,
		ErrPasswordMismatch.Error():        ErrPasswordMismatch,
		ErrAuthTokenInvalid.Error():        ErrAuthTokenInvalid,
		ErrUnexpectedSigningMethod.Error(): ErrUnexpectedSigningMethod,
		ErrAuthTokenExpired.Error():        ErrAuthTokenExpired,
		ErrAuthorIDMismatch.Error():        ErrAuthorIDMismatch,
		ErrBindRequest.Error():             ErrBindRequest,
		ErrValidateRequest.Error():         ErrValidateRequest,
		ErrDBNotFound.Error():              ErrDBNotFound,
		ErrBadJson.Error():                 ErrBadJson,
		ErrPassword.Error():                ErrPassword,
		ErrSignToken.Error():               ErrSignToken,
		ErrGenerateUUID.Error():            ErrGenerateUUID,
		ErrParseAuthToken.Error():          ErrParseAuthToken,
		ErrEmailAlreadyTaken.Error():       ErrEmailAlreadyTaken,
		ErrAddYourself.Error():             ErrAddYourself,
		ErrRequestAlreadyExist.Error():     ErrRequestAlreadyExist,
		ErrAlreadyFriends.Error():          ErrAlreadyFriends,
		ErrAlreadyFollower.Error():         ErrAlreadyFollower,
		ErrSingleChat.Error():              ErrSingleChat,
		ErrDialogAlreadyExist.Error():      ErrDialogAlreadyExist,
	}
)
