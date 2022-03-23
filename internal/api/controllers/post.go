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

type PostController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *PostController) CreatePost(ctx echo.Context) error {
	request := new(dto.GetPostDataRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.PostService.CreatePost(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *PostController) DeletePost(ctx echo.Context) error {
	request := new(dto.GetPostDeleteDataRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	err := c.registry.PostService.DeletePost(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.BasicResponse{})
}

func (c *PostController) EditPost(ctx echo.Context) error {
	request := new(dto.GetPostEditDataRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.PostService.EditPost(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewPostController(log *logrus.Entry, registry *service.Registry) *PostController {
	return &PostController{log: log, registry: registry}
}
