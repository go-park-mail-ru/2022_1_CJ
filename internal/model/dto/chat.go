package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"

type CreateDialogRequest struct {
	UserID    string   `json:"user_id"`
	AuthorIDs []string `json:"author_ids"`
}

type CreateDialogResponse struct {
	DialogID string `json:"dialog_id"`
}

type SendMessageRequest struct {
	MessageInfo common.MessageInfo `json:"message_info"`
}

type SendMessageResponse struct{}

type GetDialogsRequest struct {
	UserID string `json:"user_id"`
}

type GetDialogsResponse struct {
	DialogsInfo []common.DialogInfo `json:"dialogs_info"`
}
