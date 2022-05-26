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

type PostController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *PostController) CreatePost(ctx echo.Context) error {
	request := new(dto.CreatePostRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.PostService.CreatePost(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *PostController) GetPost(ctx echo.Context) error {
	request := new(dto.GetPostRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.PostService.GetPost(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *PostController) EditPost(ctx echo.Context) error {
	request := new(dto.EditPostRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.PostService.EditPost(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *PostController) DeletePost(ctx echo.Context) error {
	request := new(dto.DeletePostRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.PostService.DeletePost(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewPostController(log *logrus.Entry, registry *service.Registry) *PostController {
	return &PostController{log: log, registry: registry}
}
