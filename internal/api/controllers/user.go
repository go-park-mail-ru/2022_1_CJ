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
	request := new(dto.GetUserDataRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.GetUserData(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetUserFeed(ctx echo.Context) error {
	request := new(dto.GetUserFeedRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.UserService.GetUserFeed(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) SendRequest(ctx echo.Context) error {
	request := new(dto.ReqSendRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	request.PersonID = ctx.Param("person_id")

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.SendRequest(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) AcceptRequest(ctx echo.Context) error {
	request := new(dto.AcceptRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	request.PersonID = ctx.Param("person_id")

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.AcceptRequest(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) DeleteFriend(ctx echo.Context) error {
	request := new(dto.DeleteFriendRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	request.ExFriendID = ctx.Param("ex_friend_id")

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.DeleteFriend(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
