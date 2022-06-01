package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/cl"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (svc *APIService) AuthMiddlewareMicro(rep cl.AuthRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookieAuth, err := ctx.Cookie(constants.CookieKeyAuthToken)
			if err != nil {
				return constants.ErrMissingAuthCookie
			}

			newToken, UserID, code, err := rep.Check(cookieAuth.Value)
			if err != nil {
				return err
			}
			if len(UserID) != 0 {
				ctx.Request().Header.Set(constants.HeaderKeyUserID, UserID)
			}

			if code {
				ctx.SetCookie(utils.CreateHTTPOnlyCookie(constants.CookieKeyAuthToken, newToken, viper.GetInt64(constants.ViperJWTTTLKey)))
			}

			return next(ctx)
		}
	}
}

func (svc *APIService) OAuthTelegramMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			queryParams := ctx.Request().URL.Query()
			kvs := []string{}
			hash := ""
			for k, v := range queryParams {
				if k == "hash" {
					hash = v[0]
					continue
				}
				kvs = append(kvs, k+"="+v[0])
			}
			sort.Strings(kvs)

			var dataCheckString = ""
			for _, s := range kvs {
				if dataCheckString != "" {
					dataCheckString += "\n"
				}
				dataCheckString += s
			}

			sha256hash := sha256.New()

			telegramToken := viper.GetString("service.telegram_token")
			_, _ = io.WriteString(sha256hash, telegramToken)

			hmachash := hmac.New(sha256.New, sha256hash.Sum(nil))
			_, _ = io.WriteString(hmachash, dataCheckString)

			if hash != hex.EncodeToString(hmachash.Sum(nil)) {
				return constants.ErrHashInvalid
			}

			return next(ctx)
		}
	}
}

func (svc *APIService) CSRFMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookieCSRF, err := ctx.Cookie(constants.CookieKeyCSRFToken)
			if err != nil || len(cookieCSRF.Value) == 0 {
				return constants.ErrMissingCSRFCookie
			}
			tokenCSRF := ctx.QueryParam(constants.CookieKeyCSRFToken)

			if tokenCSRF != cookieCSRF.Value {
				svc.log.Errorf("Cookie token: %s; Query token: %s", cookieCSRF.Value, tokenCSRF)
				return constants.ErrCSRFTokenWrong
			}

			newTokenCSRF, err := utils.RefreshIfNeededCSRFToken(tokenCSRF, ctx.Request().Header.Get(constants.HeaderKeyUserID))
			if err != nil {
				return err
			}

			if err == nil && len(newTokenCSRF) != 0 {
				ctx.SetCookie(utils.CreateCookie(constants.CookieKeyCSRFToken, newTokenCSRF, viper.GetInt64(constants.ViperCSRFTTLKey)))
			}

			return next(ctx)
		}
	}
}

func (svc *APIService) XRequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			xRequestID := ctx.Request().Header.Get(constants.HeaderKeyRequestID)
			if len(xRequestID) == 0 {
				xRequestID, err := core.GenUUID()
				if err != nil {
					return err
				}
				ctx.Request().Header.Set(constants.HeaderKeyRequestID, xRequestID)
			}
			return next(ctx)
		}
	}
}

func (svc *APIService) LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			res := ctx.Response()

			var bodyBytes []byte
			if svc.debug {
				bodyBytes, err = ioutil.ReadAll(req.Body)
				if err != nil {
					ctx.Error(err)
				}
				req.Body.Close()
				req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			log := svc.log.WithFields(logrus.Fields{
				"x_request_id": ctx.Request().Header.Get(constants.HeaderKeyRequestID),
				"user_agent":   req.UserAgent(),
				"host":         req.Host,
				"uri":          req.RequestURI,
				"http_method":  req.Method,
				"user_ip":      getIP(req),
			})

			userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
			if len(userID) != 0 {
				log = log.WithFields(logrus.Fields{
					"user_id": userID,
				})
			}

			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()
			if res.Status >= 400 && svc.debug {
				if len(bodyBytes) > 4096 {
					bodyBytes = bodyBytes[:4096]
				}
				log = log.WithFields(logrus.Fields{
					"body": string(bodyBytes),
				})
			}

			log = log.WithFields(logrus.Fields{
				"execution_time": stop.Sub(start).String(),
				"status":         res.Status,
			})

			if res.Status >= 400 {
				log.Infof("[error]: %v", err)
			} else {
				log.Info("[success]")
			}

			return nil
		}
	}
}

func (svc *APIService) AccessLogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			res := ctx.Response()
			id, _ := uuid.NewV4()

			start := time.Now()
			ctx.Set("REQUEST_ID", id.String())

			svc.log.Infof("ID: %s; URL: %s; METHOD: %s", id.String(), ctx.Request().URL.Path, ctx.Request().Method)

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			responseTime := time.Since(start)
			svc.log.Infof("ID: %s; TIME FOR ANSWER: %s", id.String(), responseTime)

			status := res.Status
			path := ctx.Request().URL.Path
			method := ctx.Request().Method

			svc.metrics.Hits.WithLabelValues(strconv.Itoa(status), path, method).Inc()
			svc.metrics.Duration.WithLabelValues(strconv.Itoa(status), path, method).Observe(responseTime.Seconds())

			return nil
		}
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	for _, ip := range strings.Split(r.Header.Get("X-FORWARDED-FOR"), ",") {
		if netIP := net.ParseIP(ip); netIP != nil {
			return ip
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	return ""
}
