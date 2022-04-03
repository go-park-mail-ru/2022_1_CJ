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
	UserID := ctx.Param("user_id")

	response, err := c.registry.UserService.GetUserData(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetUserFeed(ctx echo.Context) error {
	UserID := ctx.Param("user_id")

	response, err := c.registry.UserService.GetUserFeed(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
