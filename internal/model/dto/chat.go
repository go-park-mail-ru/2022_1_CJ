package dto

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

// Message for chat for wb
type Message struct {
	ID        string `json:"_id"`
	DialogID  string `json:"dialog_id"`
	Event     string `json:"event"`
	AuthorID  string `json:"author_id"`
	DestinID  string `json:"dst,omitempty"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

// Message for chat for giving
type MessageInfo struct {
	AuthorID  string        `json:"author_id"`
	Body      string        `json:"body"`
	IsRead    []core.IsRead `json:"is_read,omitempty"`
	CreatedAt int64         `json:"created_at"`
}

type Dialog struct {
	DialogID     string   `json:"dialog_id"`
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
	NonRead      int64    `json:"non_read"`
}

type CreateDialogRequest struct {
	UserID    string   `json:"user_id"`
	Name      string   `json:"name"`
	AuthorIDs []string `json:"author_ids"`
}

type CreateDialogResponse struct {
	DialogID string `json:"dialog_id"`
}

type SendMessageRequest struct {
	Message Message `json:"message"`
}

type SendMessageResponse struct{}

type ReadMessageRequest struct {
	Message Message `json:"message"`
}

type ReadMessageResponse struct{}

type GetDialogsRequest struct { //
	UserID string `query:"user_id"`
	Limit  int64  `query:"limit,omitempty"`
	Page   int64  `query:"page,omitempty"`
}

type GetDialogsResponse struct {
	Dialogs     []Dialog `json:"dialogs"`
	Total       int64    `json:"total"`
	AmountPages int64    `json:"amount_pages"`
}

type GetDialogRequest struct { //
	UserID   string
	DialogID string `query:"dialog_id"`
	Limit    int64  `query:"limit,omitempty"`
	Page     int64  `query:"page,omitempty"`
}

type GetDialogResponse struct {
	Dialog      Dialog        `json:"dialog"`
	Messages    []MessageInfo `json:"messages"`
	Total       int64         `json:"total"`
	AmountPages int64         `json:"amount_pages"`
}

type CheckDialogRequest struct {
	UserID   string `json:"user_id"`
	DialogID string `json:"dialog_id"`
}
