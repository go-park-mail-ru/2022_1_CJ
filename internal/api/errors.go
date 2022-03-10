package api

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/labstack/echo"
)

func (svc *APIService) httpErrorHandler(err error, c echo.Context) {
	e := err
	msg := err.Error()
	for e != nil {
		if ce, ok := e.(*constants.CodedError); ok {
			code := ce.Code()
			if !svc.debug {
				if code == http.StatusInternalServerError {
					msg = "internal server error"
				} else {
					msg = e.Error()
				}
			}

			c.JSON(code, dto.ErrorResponse{
				Message: msg,
				Code:    code,
			})

			return
		} else {
			e = errors.Unwrap(e)
		}
	}

	if !svc.debug {
		msg = "internal server error"
	}

	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Message: msg,
		Code:    http.StatusInternalServerError,
	})
}
