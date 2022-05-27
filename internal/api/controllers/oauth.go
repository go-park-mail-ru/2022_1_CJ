package controllers

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OAuthController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *OAuthController) AuthenticateThroughTelergam(ctx echo.Context) error {
	request := new(dto.AuthenticateThroughTelergamRequest)
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.OAuthService.AuthenticateThroughTelergam(context.Background(), request)
	if err != nil {
		return err
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: request.ID})
	if err != nil {
		return err
	}

	ctx.SetCookie(utils.CreateHTTPOnlyCookie(constants.CookieKeyAuthToken, authToken, viper.GetInt64(constants.ViperJWTTTLKey)))
	ctx.SetCookie(utils.CreateCookie(constants.CookieKeyCSRFToken, response.CSRFToken, viper.GetInt64(constants.ViperCSRFTTLKey)))

	return ctx.JSON(http.StatusOK, response)
}

func NewOAuthController(log *logrus.Entry, registry *service.Registry) *OAuthController {
	return &OAuthController{log: log, registry: registry}
}
