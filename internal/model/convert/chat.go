package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Dialog2DTO(dialog *core.Dialog, userID string) dto.Dialog {
	nonRead := int64(0)
	for _, message := range dialog.Messages {
		if message.AuthorID != userID {
			for _, id := range message.IsRead {
				if id.Participant == userID {
					if !id.IsRead {
						nonRead += 1
					} else {
						break
					}
				}
			}
		} else {
			break
		}
	}

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
		NonRead:      nonRead,
	}
}

func Message2DTO(message core.Message, userID string) dto.MessageInfo {

	if userID == message.AuthorID {
		return dto.MessageInfo{
			AuthorID:  message.AuthorID,
			Body:      message.Body,
			IsRead:    message.IsRead,
			CreatedAt: message.CreatedAt,
		}
	}
	return dto.MessageInfo{
		AuthorID:  message.AuthorID,
		Body:      message.Body,
		CreatedAt: message.CreatedAt,
	}
}

func Messages2DTO(messages []core.Message, userID string) []dto.MessageInfo {
	var result []dto.MessageInfo
	for _, message := range messages {
		result = append(result, Message2DTO(message, userID))
	}
	return result
}
