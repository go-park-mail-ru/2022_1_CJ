package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	HeaderKeyUserID       = "User-Id"
	HeaderKeyRequestID    = "X-Request-Id"
	HeaderKeyUserAuthType = "User-Auth-Type"
)

const (
	CookieKeyAuthToken = "Auth-Token"
	CookieKeyCSRFToken = "X-CSRF-Token"
)

const (
	FrameRepresentativeCheckRegisterInitToken = "3"
)

const (
	QueryDelimiter = ","
)

func GetServiceEntry() string {
	return fmt.Sprintf("%s://%s", viper.GetString("service.scheme"), viper.GetString("service.host"))
}
func GetAPIEntry() string {
	return fmt.Sprintf("%s:%s/api", GetServiceEntry(), viper.GetString("service.port"))
}
