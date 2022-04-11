package chat

import (
	"context"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.send)
			}
		}
	}
}

func (h *Hub) NewClientConnectWS(ctx context.Context, conn *websocket.Conn, userID string) {
	client := NewClient(h, conn, userID)
	client.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
