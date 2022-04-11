package chat

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	ID   string
	hub  *Hub
	conn *websocket.Conn
	log  *logrus.Entry
	db   *db.Repository
	send chan common.MessageInfo
}

func NewClient(hub *Hub, conn *websocket.Conn, log *logrus.Entry,
	repo *db.Repository, userID string) *Client {
	return &Client{ID: userID,
		hub:  hub,
		conn: conn,
		log:  log,
		db:   repo,
		send: make(chan common.MessageInfo, 256),
	}
}

// Нотификация всех о сообщении
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		message := common.MessageInfo{}
		if err := c.conn.ReadJSON(&message); err != nil {
			//switch (action.type) {
			//case "S"
			//}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//log.Printf("error: %v", err)
			}
			break
		}
		message.AuthorID = c.ID
		c.send <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			json, err := json.Marshal(message)

			if err != nil {
				return
			}
			w.Write(json)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
