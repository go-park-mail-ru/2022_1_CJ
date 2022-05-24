package controllers

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FileController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *FileController) UploadFile(ctx echo.Context) error {
	image, err := ctx.FormFile("file")
	if err != nil {
		c.log.Errorf("FormFile error: %s", err)
		return err
	}

	url, err := c.registry.StaticService.UploadFile(image)
	if err != nil {
		c.log.Errorf("Upload error: %s", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.UploadFileResponse{URL: url})
}

func (c *FileController) GetFile(ctx echo.Context) error {
	request := new(dto.GetFileRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	return ctx.JSON(http.StatusOK, ctx.File("/opt/files"+request.URL))
}

func NewFileController(log *logrus.Entry, registry *service.Registry) *FileController {
	return &FileController{log: log, registry: registry}
}
