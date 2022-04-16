package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Dialog2DTO(dialog *core.Dialog, userID string) dto.Dialog {
	participants := make([]string, 0)
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
