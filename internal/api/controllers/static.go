package controllers

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type StaticController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *StaticController) UploadImage(ctx echo.Context) error {
	image, err := ctx.FormFile("image")
	if err != nil {
		return err
	}

	url, err := c.registry.StaticService.UploadImage(context.Background(), image)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.UploadImageResponse{URL: url})
}

func NewStaticController(log *logrus.Entry, registry *service.Registry) *StaticController {
	return &StaticController{log: log, registry: registry}
}
