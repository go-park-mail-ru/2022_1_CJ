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

type FriendsController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *FriendsController) SendFriendRequest(ctx echo.Context) error {
	request := new(dto.SendFriendRequestRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.SendFriendRequest(context.Background(), request, UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) AcceptFriendRequest(ctx echo.Context) error {
	request := new(dto.AcceptFriendRequestRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.FriendsService.AcceptFriendRequest(context.Background(), request, UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetFriendsByUserID(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.GetFriendsByUserID(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetOutcomingRequests(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.GetOutcomingRequests(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetIncomingRequests(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.FriendsService.GetIncomingRequests(context.Background(), UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) DeleteFriend(ctx echo.Context) error {
	request := new(dto.DeleteFriendRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	// Можно обернуть ctx в Context(ctx), чтобы передавать UserID в котексте го
	response, err := c.registry.FriendsService.DeleteFriend(context.Background(), request, UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewFriendsController(log *logrus.Entry, registry *service.Registry) *FriendsController {
	return &FriendsController{log: log, registry: registry}
}
