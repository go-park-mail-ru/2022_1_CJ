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
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	response, err := c.registry.CommunityService.GetCommunityPosts(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetUserCommunities(ctx echo.Context) error {
	request := new(dto.GetUserCommunitiesRequest)
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

	response, err := c.registry.CommunityService.GetUserCommunities(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetUserManageCommunities(ctx echo.Context) error {
	request := new(dto.GetUserManageCommunitiesRequest)
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

	response, err := c.registry.CommunityService.GetUserManageCommunities(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetCommunities(ctx echo.Context) error {
	request := new(dto.GetCommunitiesRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	response, err := c.registry.CommunityService.GetCommunities(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) CreatePostCommunity(ctx echo.Context) error {
	request := new(dto.CreatePostCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.CreatePostCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) JoinCommunity(ctx echo.Context) error {
	request := new(dto.JoinCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.JoinCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) SearchCommunities(ctx echo.Context) error {
	request := new(dto.SearchCommunitiesRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	response, err := c.registry.CommunityService.SearchCommunities(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) UpdatePhotoCommunity(ctx echo.Context) error {
	image, err := ctx.FormFile("photo")
	if err != nil {
		return err
	}

	url, err := c.registry.StaticService.UploadImage(context.Background(), image)
	if err != nil {
		return err
	}

	request := new(dto.UpdatePhotoCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	response, err := c.registry.CommunityService.UpdatePhoto(context.Background(), request, url, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) LeaveCommunity(ctx echo.Context) error {
	request := new(dto.LeaveCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.LeaveCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetFollowers(ctx echo.Context) error {
	request := new(dto.GetFollowersRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	response, err := c.registry.CommunityService.GetFollowers(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) GetMutualFriends(ctx echo.Context) error {
	request := new(dto.GetMutualFriendsRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 10
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	response, err := c.registry.CommunityService.GetMutualFriends(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) DeletePostCommunity(ctx echo.Context) error {
	request := new(dto.DeletePostCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.DeletePostCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *CommunityController) EditPostCommunity(ctx echo.Context) error {
	request := new(dto.EditPostCommunityRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)

	response, err := c.registry.CommunityService.EditPostCommunity(context.Background(), request, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

func NewCommunityController(log *logrus.Entry, registry *service.Registry) *CommunityController {
	return &CommunityController{log: log, registry: registry}
}
