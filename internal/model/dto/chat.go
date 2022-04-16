package dto

type Message struct {
	DialogID string `json:"dialog_id"`
	Event    string `json:"event"`
	AuthorID string `json:"author_id"`
	DestinID string `json:"dst,omitempty"`
	Body     string `json:"body"`
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
	MessageInfo Message `json:"message"`
}

type SendMessageResponse struct{}

type GetDialogsRequest struct {
	UserID string `json:"user_id"`
}

type GetDialogsResponse struct {
	Dialogs []Dialog `json:"dialogs"`
}
