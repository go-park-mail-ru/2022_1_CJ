package controllers

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *UserController) GetUserData(ctx echo.Context) error {
	request := new(dto.GetUserRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.GetUserData(context.Background(), request.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetUserPosts(ctx echo.Context) error {
	request := new(dto.GetUserPostsRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.GetUserPosts(context.Background(), request.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetFeed(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.UserService.GetFeed(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
