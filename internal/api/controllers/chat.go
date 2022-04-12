package controllers

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core/chat"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ChatController struct {
	hub      *chat.Hub
	log      *logrus.Entry
	registry *service.Registry
	db       *db.Repository
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *ChatController) GetChats(ctx echo.Context) error {
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

func (c *ChatController) CreateChat(ctx echo.Context) error {
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
	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		c.log.Errorf("Error upgarde: %s", err)
		return err
	}
	defer conn.Close()
	UserID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	c.hub.NewClientConnectWS(context.Background(), conn, c.log, c.db, UserID)
	return nil
}

func NewChatController(hub *chat.Hub, log *logrus.Entry, repo *db.Repository, registry *service.Registry) *ChatController {
	return &ChatController{hub: hub, log: log, db: repo, registry: registry}
}
