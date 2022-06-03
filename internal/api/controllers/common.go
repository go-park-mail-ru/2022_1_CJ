package controllers

import (
	"strconv"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/labstack/echo/v4"
)

func parsePaginationParametersQuery(ctx echo.Context) (core.PaginationParameters, error) {
	params := core.PaginationParameters{}

	if limit, err := strconv.ParseInt(ctx.QueryParam("limit"), 10, 64); err == nil && limit >= -1 {
		params.Limit = limit
	} else {
		params.Limit = 10
	}

	if page, err := strconv.ParseInt(ctx.QueryParam("page"), 10, 64); err == nil {
		params.Page = page
	} else {
		params.Page = 1
	}

	return params, nil
}
