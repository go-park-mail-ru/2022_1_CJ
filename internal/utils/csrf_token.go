package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

func GenerateCSRFToken(userID string) (string, error) {
	h := hmac.New(sha256.New, []byte(constants.ViperCSRFSecretKey))

	t := time.Second * time.Duration(viper.GetInt64(constants.ViperCSRFTTLKey))
	timeNow := time.Now().Add(t).Unix()

	data := fmt.Sprintf("%s:%d", userID, timeNow)
	h.Write([]byte(data))

	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(timeNow, 10)

	return token, nil
}

func RefreshIfNeededCSRFToken(token string, userID string) (string, error) {
	tokenData := strings.Split(token, ":")

	if len(tokenData) != 2 {
		return "", constants.ErrCSRFTokenWrong
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return "", constants.ErrCSRFTokenWrong
	}

	if tokenExp > time.Now().Unix()+viper.GetInt64(constants.ViperCSRFTTLKey)/2 {
		return "", nil
	}

	return GenerateCSRFToken(userID)
}
