package constants

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	WriteWait      = 600 * time.Second
	PongWait       = 600 * time.Second
	PingPeriod     = PongWait * 9 / 10
	MaxMessageSize = 1024 * 1024 * 1024

	JoinChat    = "join"
	LeaveChat   = "leave"
	JoinedChat  = "joined"
	LeftChat    = "left"
	SendChat    = "send"
	SendFile    = "send_file"
	SendSticker = "send_sticker"
	ReadChat    = "read"
	Empty       = ""

	ErrChat           = "error"
	ErrChatDoNotExist = "room does not exit"
	ErrRequest        = "bad request"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
