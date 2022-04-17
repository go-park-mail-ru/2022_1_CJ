package chat

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// The Conn type represents a single client.
type Conn struct {
	sync.Mutex
	Socket  *websocket.Conn
	ID      string
	Send    chan dto.Message
	Dialogs map[string]string
	reg     *service.Registry
	log     *logrus.Entry
}

var (
	// Stores all Conn types by their uuid.
	ConnManager = struct {
		sync.Mutex
		Conns map[string]*Conn
	}{
		Conns: make(map[string]*Conn),
	}
)

// Handles incoming, error free messages.
func HandleData(c *Conn, msg *dto.Message) {
	switch msg.Event {
	case constants.JoinChat:
		c.Join(msg.DialogID)
	case constants.LeaveChat:
		c.Leave(msg.DialogID)
	case constants.JoinedChat:
		c.Emit(msg)
	case constants.LeftChat:
		c.Emit(msg)
		DialogManager.Lock()
		room, ok := DialogManager.Rooms[msg.DialogID]
		DialogManager.Unlock()
		if ok == false {
			break
		}
		room.Lock()
		delete(room.Members, c.ID)
		members := len(room.Members)
		room.Unlock()
		if members == 0 {
			room.Stop()
		}
	case constants.ReadChat:
		if msg.DestinID != constants.Empty {
			DialogManager.Lock()
			room, rok := DialogManager.Rooms[msg.DialogID]
			DialogManager.Unlock()
			if rok == false {
				break
			}
			room.Lock()
			id, mok := room.Members[msg.DestinID]
			room.Unlock()
			if mok == false {
				break
			}
			ConnManager.Lock()
			dst, cok := ConnManager.Conns[id]
			ConnManager.Unlock()
			if cok == false {
				break
			}
			if msg.Body != constants.Empty {
				_, err := c.reg.ChatService.ReadMessage(context.Background(), &dto.ReadMessageRequest{Message: *msg})
				if err != nil {
					return
				}
				dst.Send <- *msg
			}
		}
	case constants.SendChat:
		msgID, err := core.GenUUID()
		if err != nil {
			break
		}
		msg.ID = msgID
		msg.CreatedAt = time.Now().Unix()

		_, err = c.reg.ChatService.SendMessage(context.Background(), &dto.SendMessageRequest{Message: *msg})
		if err != nil {
			return
		}
		c.Emit(msg)
	default:
		return
	}
}
func (c *Conn) readPump() {
	defer func() {
		c.Lock()
		for name := range c.Dialogs {
			c.Unlock()
			DialogManager.Lock()
			room, ok := DialogManager.Rooms[name]
			DialogManager.Unlock()
			if ok == true {
				room.Leave(c)
			}
			c.Lock()
		}
		c.Unlock()
		c.Socket.Close()
	}()
	c.Socket.SetReadLimit(constants.MaxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(constants.PongWait))
		return nil
	})
	for {
		var data dto.Message
		err := c.Socket.ReadJSON(&data)
		data.AuthorID = c.ID

		if err != nil {
			if _, wok := err.(*websocket.CloseError); wok == false {
				break
			}
			c.Lock()
			for name := range c.Dialogs {
				c.Unlock()
				DialogManager.Lock()
				room, rok := DialogManager.Rooms[name]
				DialogManager.Unlock()
				if rok == false {
					c.Lock()
					continue
				}
				room.Emit(c, ConstructMessage(name, constants.LeftChat, c.ID, constants.Empty, c.ID))
				room.Lock()
				delete(room.Members, c.ID)
				members := len(room.Members)
				room.Unlock()
				if members == 0 {
					room.Stop()
				}
				c.Lock()
			}
			c.Unlock()
			break
		}
		HandleData(c, &data)
	}
}

func (c *Conn) write(mt int, payload []byte) error {
	c.Socket.SetWriteDeadline(time.Now().Add(constants.WriteWait))
	return c.Socket.WriteMessage(mt, payload)
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if ok == false {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Socket.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Adds the Conn to a Dialog. If the Dialog does not exist, it is created.
func (c *Conn) Join(name string) {
	DialogManager.Lock()
	room, ok := DialogManager.Rooms[name]
	DialogManager.Unlock()
	if ok == false {
		room = NewRoom(name)
	}
	c.Lock()
	c.Dialogs[name] = name
	c.Unlock()
	room.Join(c)
}

// Removes the Conn from a Dialog.
func (c *Conn) Leave(name string) {
	DialogManager.Lock()
	room, rok := DialogManager.Rooms[name]
	DialogManager.Unlock()
	if rok == false {
		return
	}
	c.Lock()
	_, cok := c.Dialogs[name]
	c.Unlock()
	if cok == false {
		return
	}
	c.Lock()
	delete(c.Dialogs, name)
	c.Unlock()
	room.Leave(c)
}

// Broadcasts a Message to all members of a Dialog.
func (c *Conn) Emit(msg *dto.Message) {
	DialogManager.Lock()
	room, ok := DialogManager.Rooms[msg.DialogID]
	DialogManager.Unlock()
	if ok == true {
		room.Emit(c, msg)
	}
}

// Upgrades an HTTP connection and creates a new Conn type.
func NewConnection(ctx *echo.Context, log *logrus.Entry, registry *service.Registry, userID string) *Conn {
	socket, err := constants.Upgrader.Upgrade((*ctx).Response(), (*ctx).Request(), nil)
	if err != nil {
		return nil
	}
	c := &Conn{
		Socket:  socket,
		ID:      userID,
		Send:    make(chan dto.Message),
		Dialogs: make(map[string]string),
		log:     log,
		reg:     registry,
	}
	ConnManager.Lock()
	ConnManager.Conns[c.ID] = c
	ConnManager.Unlock()
	return c
}

// Calls NewConnection, starts the returned Conn's writer, joins the root room, and finally starts the Conn's reader.
func SocketHandler(ctx *echo.Context, log *logrus.Entry, registry *service.Registry, userID string) error {
	c := NewConnection(ctx, log, registry, userID)
	if c != nil {
		go c.writePump()
		go c.readPump()
		c.log.Infof("new user: %s", c.ID)
	}
	return nil
}
