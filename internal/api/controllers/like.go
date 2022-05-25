package controllers

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type LikeController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *LikeController) IncreaseLike(ctx echo.Context) error {
	request := new(dto.IncreaseLikeRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.LikeService.IncreaseLike(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *LikeController) ReduceLike(ctx echo.Context) error {
	request := new(dto.ReduceLikeRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.LikeService.ReduceLike(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *LikeController) GetLikePost(ctx echo.Context) error {
	request := new(dto.GetLikePostRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.LikeService.GetLikePost(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *LikeController) GetLikePhoto(ctx echo.Context) error {
	request := new(dto.GetLikePhotoRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.LikeService.GetLikePhoto(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewLikeController(log *logrus.Entry, registry *service.Registry) *LikeController {
	return &LikeController{log: log, registry: registry}
}
