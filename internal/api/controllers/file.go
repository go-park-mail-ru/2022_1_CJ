package controllers

import (
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type FileController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *FileController) UploadFile(ctx echo.Context) error {
	image, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	url, err := c.registry.StaticService.UploadFile(image)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.UploadFileResponse{URL: url})
}

func (c *FileController) GetFile(ctx echo.Context) error {
	request := new(dto.GetFileRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return ctx.Inline("/opt/files"+request.URL, request.URL)
}

func NewFileController(log *logrus.Entry, registry *service.Registry) *FileController {
	return &FileController{log: log, registry: registry}
}
