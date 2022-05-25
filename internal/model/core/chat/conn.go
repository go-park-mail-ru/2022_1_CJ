package chat

import (
	"context"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

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
	ctx     echo.Context
}

var (
	// ConnManager Stores all Conn types by their uuid.
	ConnManager = struct {
		sync.Mutex
		Conns map[string]*Conn
	}{
		Conns: make(map[string]*Conn),
	}
)

// HandleData Handles incoming, error free messages.
func HandleData(c *Conn, msg *dto.Message) {
	if c.IsDialogExist(msg.DialogID) {
		switch msg.Event {
		case constants.JoinChat:
			c.Join(msg.DialogID)
		case constants.LeaveChat:
			c.Leave(msg.DialogID)
		case constants.JoinedChat:
			c.Emit(msg)
		case constants.LeftChat:
			c.LeftChat(msg)
		case constants.ReadChat:
			c.ReadMessage(msg)
		case constants.SendChat:
			c.SendMessage(msg)
		case constants.SendFile:
			c.SendFile(msg)
		case constants.SendSticker:
			c.SendSticker(msg)
		default:
			c.Send <- *ConstructMessage(msg.DialogID, constants.ErrChat, c.ID, constants.Empty, constants.ErrRequest)
		}
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
		_ = c.Socket.Close()
	}()
	c.Socket.SetReadLimit(constants.MaxMessageSize)
	_ = c.Socket.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.Socket.SetPongHandler(func(string) error {
		_ = c.Socket.SetReadDeadline(time.Now().Add(constants.PongWait))
		return nil
	})
	for {
		data := new(dto.Message)
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
		HandleData(c, data)
	}
}

func (c *Conn) write(mt int, payload []byte) error {
	_ = c.Socket.SetWriteDeadline(time.Now().Add(constants.WriteWait))
	return c.Socket.WriteMessage(mt, payload)
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Socket.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.write(websocket.CloseMessage, []byte{})
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

// Join Adds the Conn to a Dialog. If the Dialog does not exist, it is created.
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

// SendMessage ...
func (c *Conn) SendMessage(msg *dto.Message) {
	msgID, err := core.GenUUID()
	if err != nil {
		return
	}
	msg.ID = msgID
	msg.CreatedAt = time.Now().Unix()

	_, err = c.reg.ChatService.SendMessage(context.Background(), &dto.SendMessageRequest{Message: *msg})
	if err != nil {
		c.log.Errorf("don't send message: %s", err)
		return
	}
	c.log.Info("send message")
	c.Emit(msg)
}

func (c *Conn) SendFile(msg *dto.Message) {
	ext := c.ctx.FormValue("ext")
	if ext == "" {
		c.log.Error("Extension is empty")
		return
	}

	msgID, err := core.GenUUID()
	if err != nil {
		return
	}
	msg.ID = msgID
	msg.CreatedAt = time.Now().Unix()

	msg.Body, err = c.reg.StaticService.UploadFileChat(msgID, ext, msg.Body)
	if err != nil {
		c.log.Errorf("UploadFileChat error: %s", err)
		return
	}

	_, err = c.reg.ChatService.SendMessage(context.Background(), &dto.SendMessageRequest{Message: *msg})
	if err != nil {
		c.log.Errorf("don't send message: %s", err)
		return
	}
	c.log.Infof("send message")
	c.Emit(msg)
}

// SendMessage ...
func (c *Conn) SendSticker(msg *dto.Message) {
	msgID, err := core.GenUUID()
	if err != nil {
		return
	}
	msg.ID = msgID
	msg.CreatedAt = time.Now().Unix()

	_, err = c.reg.ChatService.SendMessage(context.Background(), &dto.SendMessageRequest{Message: *msg})
	if err != nil {
		c.log.Errorf("don't send message: %s", err)
		return
	}
	c.log.Info("send message")
	c.Emit(msg)
}

// LeftChat ...
func (c *Conn) LeftChat(msg *dto.Message) {
	c.Emit(msg)
	DialogManager.Lock()
	room, ok := DialogManager.Rooms[msg.DialogID]
	DialogManager.Unlock()
	if ok == false {
		return
	}
	room.Lock()
	delete(room.Members, c.ID)
	members := len(room.Members)
	room.Unlock()
	if members == 0 {
		room.Stop()
	}
}

// ReadMessage ...
func (c *Conn) ReadMessage(msg *dto.Message) {
	if msg.DestinID != constants.Empty {
		DialogManager.Lock()
		room, rok := DialogManager.Rooms[msg.DialogID]
		DialogManager.Unlock()
		if rok == false {
			return
		}
		room.Lock()
		id, mok := room.Members[msg.DestinID]
		room.Unlock()
		if mok == false {
			return
		}
		ConnManager.Lock()
		dst, cok := ConnManager.Conns[id]
		ConnManager.Unlock()
		if cok == false {
			return
		}
		if msg.Body != constants.Empty {

			_, err := c.reg.ChatService.ReadMessage(context.Background(), &dto.ReadMessageRequest{Message: *msg})
			if err != nil {
				c.log.Errorf("don't read message in db")
				return
			}
			c.log.Infof("read message")
			dst.Send <- *msg
		}
	}
}

// Leave Removes the Conn from a Dialog.
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

// IsDialogExist CheckDialog ...
func (c *Conn) IsDialogExist(name string) bool {

	err := c.reg.ChatService.CheckDialog(context.Background(), &dto.CheckDialogRequest{UserID: c.ID, DialogID: name})
	if err != nil {
		c.Send <- *ConstructMessage(name, constants.ErrChat, c.ID, constants.Empty, constants.ErrChatDoNotExist)
		c.log.Debug("room doesn't exist")
		return false
	}
	return true
}

// Emit Broadcasts a Message to all members of a Dialog.
func (c *Conn) Emit(msg *dto.Message) {
	DialogManager.Lock()
	room, ok := DialogManager.Rooms[msg.DialogID]
	DialogManager.Unlock()
	if ok == true {
		room.Emit(c, msg)
	}
}

// NewConnection Upgrades an HTTP connection and creates a new Conn type.
func NewConnection(ctx *echo.Context, log *logrus.Entry, registry *service.Registry, userID string) (*Conn, error) {
	socket, err := constants.Upgrader.Upgrade((*ctx).Response(), (*ctx).Request(), nil)
	if err != nil {
		return nil, err
	}
	conn := &Conn{
		Socket:  socket,
		ID:      userID,
		Send:    make(chan dto.Message),
		Dialogs: make(map[string]string),
		log:     log,
		reg:     registry,
		ctx:     *ctx,
	}
	ConnManager.Lock()
	ConnManager.Conns[conn.ID] = conn
	ConnManager.Unlock()
	return conn, nil
}

// SocketHandler Calls NewConnection, starts the returned Conn's writer, joins the root room, and finally starts the Conn's reader.
func SocketHandler(ctx *echo.Context, log *logrus.Entry, registry *service.Registry, userID string) error {
	conn, err := NewConnection(ctx, log, registry, userID)
	if err != nil {
		return err
	}

	go conn.writePump()
	go conn.readPump()
	conn.log.Infof("new user: %s", conn.ID)

	return nil
}
