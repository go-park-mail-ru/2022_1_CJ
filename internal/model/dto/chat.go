package dto

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
	AuthorID  string `json:"author_id"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt int64  `json:"created_at"`
}

type Dialog struct {
	DialogID     string   `json:"dialog_id"`
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

type CreateDialogRequest struct {
	UserID    string   `json:"user_id"`
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

type GetDialogsRequest struct {
	UserID string `json:"user_id"`
}

type GetDialogsResponse struct {
	Dialogs []Dialog `json:"dialogs"`
}

type GetDialogRequest struct {
	UserID   string `json:"user_id"`
	DialogID string `json:"dialog_id"`
}

type GetDialogResponse struct {
	Dialog   Dialog        `json:"dialog"`
	Messages []MessageInfo `json:"messages"`
}
