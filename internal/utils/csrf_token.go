package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/spf13/viper"
	"strconv"
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
