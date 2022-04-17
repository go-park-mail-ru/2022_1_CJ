package controllers

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core/chat"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ChatController struct {
	log      *logrus.Entry
	registry *service.Registry
	db       *db.Repository
}

func (c *ChatController) GetDialogs(ctx echo.Context) error {
	request := new(dto.GetDialogsRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.ChatService.GetDialogs(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ChatController) GetDialog(ctx echo.Context) error {
	request := new(dto.GetDialogRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.ChatService.GetDialog(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ChatController) CreateDialog(ctx echo.Context) error {
	request := new(dto.CreateDialogRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.ChatService.CreateDialog(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ChatController) WsHandler(ctx echo.Context) error {
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	return chat.SocketHandler(&ctx, c.log, c.registry, userID)
}

func NewChatController(log *logrus.Entry, repo *db.Repository, registry *service.Registry) *ChatController {
	return &ChatController{log: log, db: repo, registry: registry}
}
