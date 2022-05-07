package controllers

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type FriendsController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *FriendsController) SendRequest(ctx echo.Context) error {
	request := new(dto.SendFriendRequestRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.SendRequest(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) RevokeRequest(ctx echo.Context) error {
	request := new(dto.RevokeFriendRequestRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.RevokeRequest(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) AcceptRequest(ctx echo.Context) error {
	request := new(dto.AcceptFriendRequestRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.AcceptRequest(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetFriends(ctx echo.Context) error {
	request := new(dto.GetFriendsRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.GetFriends(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) DeleteFriend(ctx echo.Context) error {
	request := new(dto.DeleteFriendRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.DeleteFriend(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetIncomingRequests(ctx echo.Context) error {
	request := new(dto.GetIncomingRequestsRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.GetIncomingRequests(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *FriendsController) GetOutcomingRequests(ctx echo.Context) error {
	request := new(dto.GetOutcomingRequestsRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.FriendsService.GetOutcomingRequests(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewFriendsController(log *logrus.Entry, registry *service.Registry) *FriendsController {
	return &FriendsController{log: log, registry: registry}
}
