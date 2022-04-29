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
	data := fmt.Sprintf("%s:%d", userID, viper.GetInt64(constants.ViperCSRFTTLKey))
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(viper.GetInt64(constants.ViperCSRFTTLKey), 10)
	//fmt.Print(strconv.FormatInt(viper.GetInt64(constants.ViperCSRFTTLKey, 10)))
	return token, nil
}

func CheckCSRFToken(userID string, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	h := hmac.New(sha256.New, []byte(constants.ViperCSRFSecretKey))
	data := fmt.Sprintf("%s:%d", userID, viper.GetInt64(constants.ViperCSRFTTLKey))
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("cand hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
