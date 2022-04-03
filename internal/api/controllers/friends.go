package controllers

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FriendsController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *FriendsController) SendRequest(ctx echo.Context) error {
	PersonID := ctx.Param("person_id")
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.SendRequest(context.Background(), UserID, PersonID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) AcceptRequest(ctx echo.Context) error {
	request := new(dto.AcceptRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	PersonID := ctx.Param("person_id")
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.AcceptRequest(context.Background(), request, UserID, PersonID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) DeleteFriend(ctx echo.Context) error {
	ExFriendID := ctx.Param("ex_friend_id")
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	// Можно обернуть ctx в Context(ctx), чтобы передавать UserID в котексте го
	response, err := c.registry.FriendsService.DeleteFriend(context.Background(), UserID, ExFriendID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetFriends(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.GetFriends(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetRequests(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.GetRequests(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewFriendsController(log *logrus.Entry, registry *service.Registry) *FriendsController {
	return &FriendsController{log: log, registry: registry}
}
