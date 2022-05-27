package controllers

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/cl"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AuthController struct {
	log      *logrus.Entry
	registry *service.Registry
	rep      cl.AuthRepository
}

func (c *AuthController) SignupUser(ctx echo.Context) error {
	request := new(dto.SignupUserRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	token, userID, err := c.rep.SignUp(request.Email, request.Password)
	if err != nil {
		return err
	}

	response, err := c.registry.AuthService.SignupUser(context.Background(), request, userID, token)
	if err != nil {
		return err
	}

	ctx.SetCookie(utils.CreateHTTPOnlyCookie(constants.CookieKeyAuthToken, token, viper.GetInt64(constants.ViperJWTTTLKey)))
	ctx.SetCookie(utils.CreateCookie(constants.CookieKeyCSRFToken, response.CSRFToken, viper.GetInt64(constants.ViperCSRFTTLKey)))

	return ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) LoginUser(ctx echo.Context) error {
	request := new(dto.LoginUserRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	token, userID, err := c.rep.Login(request.Email, request.Password)
	if err != nil {
		return err
	}
	response, err := c.registry.AuthService.LoginUser(context.Background(), userID, token)
	if err != nil {
		return err
	}

	ctx.SetCookie(utils.CreateHTTPOnlyCookie(constants.CookieKeyAuthToken, token, viper.GetInt64(constants.ViperJWTTTLKey)))
	ctx.SetCookie(utils.CreateCookie(constants.CookieKeyCSRFToken, response.CSRFToken, viper.GetInt64(constants.ViperCSRFTTLKey)))

	return ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) LogoutUser(ctx echo.Context) error {
	ctx.SetCookie(utils.CreateHTTPOnlyCookie(constants.CookieKeyAuthToken, "", 0))
	ctx.SetCookie(utils.CreateCookie(constants.CookieKeyCSRFToken, "", 0))
	return ctx.JSON(http.StatusOK, &dto.BasicResponse{})
}

func NewAuthController(log *logrus.Entry, registry *service.Registry, rep cl.AuthRepository) *AuthController {
	return &AuthController{log: log, registry: registry, rep: rep}
}
