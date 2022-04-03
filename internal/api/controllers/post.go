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

	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.PostService.CreatePost(context.Background(), request, UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *PostController) DeletePost(ctx echo.Context) error {
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	PostID := ctx.Param("post_id")

	err := c.registry.PostService.DeletePost(context.Background(), UserID, PostID)
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

	PostID := ctx.Param("post_id")
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.PostService.EditPost(context.Background(), request, UserID, PostID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *PostController) GetPost(ctx echo.Context) error {

	PostID := ctx.Param("post_id")
	response, err := c.registry.PostService.GetPost(context.Background(), PostID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewPostController(log *logrus.Entry, registry *service.Registry) *PostController {
	return &PostController{log: log, registry: registry}
}
