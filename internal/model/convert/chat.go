package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Dialog2DTO(dialog *core.Dialog, userID string) dto.Dialog {
	var participants []string
	for _, id := range dialog.Participants {
		if id != userID {
			participants = append(participants, id)
		}
	}
	return dto.Dialog{
		DialogID:     dialog.ID,
		Name:         dialog.Name,
		Participants: participants,
	}
}

func Message2DTO(messages core.Message) dto.MessageInfo {
	return dto.MessageInfo{
		AuthorID:  messages.AuthorID,
		Body:      messages.Body,
		IsRead:    messages.IsRead,
		CreatedAt: messages.CreatedAt,
	}
}

func Messages2DTO(messages []core.Message) []dto.MessageInfo {
	var result []dto.MessageInfo
	for _, message := range messages {
		result = append(result, Message2DTO(message))
	}
	return result
}
