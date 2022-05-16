package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDialog2DTO(t *testing.T) {
	dialogCore := &core.Dialog{ID: "123", Name: "bestChat", Participants: []string{"123", "1233"}, CreatedAt: 123}
	dialogDTO := convert.Dialog2DTO(dialogCore, "5")
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, dialogDTO, dto.Dialog{DialogID: "123", Name: "bestChat", Participants: []string{"123", "1233"}}) {
			t.Error("got : ", dialogDTO, " expected :", dto.Dialog{DialogID: "123", Name: "bestChat", Participants: []string{"123", "1233"}})
		}
	})
}

func TestMessage2DTO(t *testing.T) {
	messageCore := core.Message{ID: "123", Body: "someBody", AuthorID: "1234", CreatedAt: 123}
	messageDTO := convert.Message2DTO(messageCore, "5")
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, messageDTO, dto.MessageInfo{AuthorID: "1234", Body: "someBody", CreatedAt: 123}) {
			t.Error("got : ", messageDTO, " expected :", dto.MessageInfo{AuthorID: "1234", Body: "someBody", CreatedAt: 123})
		}
	})
}

func TestMessages2DTO(t *testing.T) {
	messagesCore := []core.Message{{ID: "123", Body: "someBody", AuthorID: "1234", CreatedAt: 123}, {ID: "1234", Body: "someBody", AuthorID: "123", CreatedAt: 1235}}
	messagesDTO := convert.Messages2DTO(messagesCore, "5")
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, messagesDTO[0], dto.MessageInfo{AuthorID: "1234", Body: "someBody", CreatedAt: 123}) {
			t.Error("got : ", messagesDTO[0], " expected :", dto.MessageInfo{AuthorID: "1234", Body: "someBody", CreatedAt: 123})
		}
		if !assert.Equal(t, messagesDTO[1], dto.MessageInfo{AuthorID: "123", Body: "someBody", CreatedAt: 1235}) {
			t.Error("got : ", messagesDTO[1], " expected :", dto.MessageInfo{AuthorID: "123", Body: "someBody", CreatedAt: 1235})
		}
	})
}
