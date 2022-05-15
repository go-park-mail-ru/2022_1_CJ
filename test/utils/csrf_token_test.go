package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestGenerateCSRFToken(t *testing.T) {
	// Test 1
	testUserID := "123"
	csrfFirst, _ := utils.GenerateCSRFToken(testUserID)
	csrfSecond, _ := utils.GenerateCSRFToken(testUserID)
	t.Run("Check cookie", func(t *testing.T) {
		if !assert.Equal(t, csrfFirst, csrfSecond) {
			t.Error("got : ", csrfFirst, " expected :", csrfSecond)
		}
	})
	// Test 2
	testNewUserID := "223"
	csrfThird, _ := utils.GenerateCSRFToken(testNewUserID)
	t.Run("Check cookie", func(t *testing.T) {
		if !assert.NotEqual(t, csrfFirst, csrfThird) {
			t.Error("got : ", csrfFirst, "not expected :", csrfSecond)
		}
	})
}

func TestRefreshIfNeededCSRFToken(t *testing.T) {
	// Test Not equal
	h := hmac.New(sha256.New, []byte(constants.ViperCSRFSecretKey))
	timeD := time.Second * time.Duration(1)
	timeNow := time.Now().Add(timeD).Unix()
	data := fmt.Sprintf("%s:%d", "123", timeNow)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(timeNow, 10)
	tokenNew, _ := utils.RefreshIfNeededCSRFToken(token, "123")
	t.Run("Check cookie", func(t *testing.T) {
		if !assert.NotEqual(t, token, tokenNew) {
			t.Error("got : ", tokenNew, "not expected :", token)
		}
	})

	// Test Not equal
	hNew := hmac.New(sha256.New, []byte(constants.ViperCSRFSecretKey))
	timeDNew := time.Second * time.Duration(604800)
	timeNowNew := time.Now().Add(timeDNew).Unix()
	dataNew := fmt.Sprintf("%s:%d", "123", timeNowNew)
	hNew.Write([]byte(dataNew))
	tokenS := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(timeNowNew, 10)
	tokenNewS, _ := utils.RefreshIfNeededCSRFToken(tokenS, "123")
	t.Run("Check cookie", func(t *testing.T) {
		if !assert.Equal(t, "", tokenNewS) {
			t.Error("got : ", tokenNewS, "expected :", "")
		}
	})
}
