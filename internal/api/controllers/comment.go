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

type CommentController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *CommentController) CreateComment(ctx echo.Context) error {
	request := new(dto.CreateCommentRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommentService.CreateComment(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) GetComments(ctx echo.Context) error {
	request := new(dto.GetCommentsRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	response, err := c.registry.CommentService.GetComments(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) EditComment(ctx echo.Context) error {
	request := new(dto.EditCommentRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.CommentService.EditComment(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) DeleteComment(ctx echo.Context) error {
	request := new(dto.DeleteCommentRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.CommentService.DeleteComment(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewCommentController(log *logrus.Entry, registry *service.Registry) *CommentController {
	return &CommentController{log: log, registry: registry}
}
