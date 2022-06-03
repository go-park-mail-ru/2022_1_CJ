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

type UserController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *UserController) GetUserData(ctx echo.Context) error {
	request := new(dto.GetUserRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.GetUserData(context.Background(), request.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetUserPosts(ctx echo.Context) error {
	request := new(dto.GetUserPostsRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	response, err := c.registry.UserService.GetUserPosts(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetFeed(ctx echo.Context) error {
	request := new(dto.GetUserFeedRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	paginationParameters, err := parsePaginationParametersQuery(ctx)
	if err != nil {
		return err
	}

	request.PaginationParameters = paginationParameters
	response, err := c.registry.UserService.GetFeed(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetProfile(ctx echo.Context) error {
	request := new(dto.GetProfileRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	if len(request.UserID) == 0 {
		request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	}

	response, err := c.registry.UserService.GetProfile(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) EditProfile(ctx echo.Context) error {
	request := new(dto.EditProfileRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.UserService.EditProfile(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) UpdatePhoto(ctx echo.Context) error {
	image, err := ctx.FormFile("photo")
	if err != nil {
		return err
	}

	url, err := c.registry.StaticService.UploadImage(context.Background(), image)
	if err != nil {
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.UserService.UpdatePhoto(context.Background(), url, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserController) SearchUsers(ctx echo.Context) error {
	request := new(dto.SearchUsersRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	response, err := c.registry.UserService.SearchUsers(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
