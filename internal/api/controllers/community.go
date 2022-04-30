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

type CommunityController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *CommunityController) CreateCommunity(ctx echo.Context) error {
	request := new(dto.CreateCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.CreateCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) DeleteCommunity(ctx echo.Context) error {
	request := new(dto.DeleteCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.DeleteCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) EditCommunity(ctx echo.Context) error {
	request := new(dto.EditCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.EditCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetCommunity(ctx echo.Context) error {
	request := new(dto.GetCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	response, err := c.registry.CommunityService.GetCommunity(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetCommunityPosts(ctx echo.Context) error {
	request := new(dto.GetCommunityPostsRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.GetCommunityPosts(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewCommunityController(log *logrus.Entry, registry *service.Registry) *CommunityController {
	return &CommunityController{log: log, registry: registry}
}
