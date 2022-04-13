package chat

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// The Conn type represents a single client.
type Conn struct {
	sync.Mutex
	Socket *websocket.Conn
	ID     string
	Send   chan *Message
	Rooms  map[string]string
	db     *db.Repository
	log    *logrus.Entry
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 1024 * 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
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
func HandleData(c *Conn, msg *Message) {
	switch msg.Event {
	case "join":
		c.Join(msg.DialogID)
	case "leave":
		c.Leave(msg.DialogID)
	case "joined":
		c.Emit(msg)
	case "left":
		c.Emit(msg)
		RoomManager.Lock()
		room, ok := RoomManager.Rooms[msg.DialogID]
		RoomManager.Unlock()
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
	default:
		if msg.DialogID != "" {
			c.log.Infof("Write in dialog")
			RoomManager.Lock()
			room, rok := RoomManager.Rooms[msg.DialogID]
			RoomManager.Unlock()
			if rok == false {
				break
			}
			room.Lock()
			_, mok := room.Members[msg.DialogID]
			room.Unlock()
			if mok == false {
				break
			}
			ConnManager.Lock()
			// Strange
			// was msg.SrcID
			dst, cok := ConnManager.Conns[c.ID]
			ConnManager.Unlock()
			if cok == false {
				break
			}
			dst.Send <- msg
		} else {
			c.Emit(msg)
		}
	}
}

func (c *Conn) readPump() {
	defer func() {
		c.Lock()
		for id := range c.Rooms {
			c.Unlock()
			RoomManager.Lock()
			room, ok := RoomManager.Rooms[id]
			RoomManager.Unlock()
			if ok == true {
				room.Leave(c)
			}
			c.Lock()
		}
		c.Unlock()
		c.Socket.Close()
	}()
	c.Socket.SetReadLimit(maxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		data := new(Message)
		err := c.Socket.ReadJSON(&data)
		c.log.Infof("readPump smth: %s", data.Event)
		if err != nil {
			if _, wok := err.(*websocket.CloseError); wok == false {
				break
			}
			c.Lock()
			for id := range c.Rooms {
				c.Unlock()
				RoomManager.Lock()
				room, rok := RoomManager.Rooms[id]
				RoomManager.Unlock()
				if rok == false {
					c.Lock()
					continue
				}
				room.Emit(c, ConstructMessage(id, "left", "", c.ID, []byte(c.ID)))
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
	c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
	//c.Socket.WriteJSON()
	return c.Socket.WriteMessage(mt, payload)
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			c.log.Infof("writePump smth: %s", msg.Event)
			if ok == false {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteJSON(msg)
			bytes, err := json.Marshal(msg)
			if err != nil {
				c.log.Infof("true")
				c.write(websocket.CloseMessage, []byte{})
				return
			} else {
				c.log.Infof("false")
			}
			if err := c.write(websocket.BinaryMessage, bytes); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Adds the Conn to a Room. If the Room does not exist, it is created.
func (c *Conn) Join(id string) {
	c.log.Infof("join conn method with id: %s", id)
	RoomManager.Lock()
	room, ok := RoomManager.Rooms[id]
	RoomManager.Unlock()
	if ok == false {
		room = NewRoom(id)
	}
	c.Lock()
	c.Rooms[id] = id
	c.Unlock()
	room.Join(c)
}

// Removes the Conn from a Room.
func (c *Conn) Leave(id string) {
	c.log.Infof("leave conn method with id: %s", id)
	RoomManager.Lock()
	room, rok := RoomManager.Rooms[id]
	RoomManager.Unlock()
	if rok == false {
		return
	}
	c.Lock()
	_, cok := c.Rooms[id]
	c.Unlock()
	if cok == false {
		return
	}
	c.Lock()
	delete(c.Rooms, id)
	c.Unlock()
	room.Leave(c)
}

// Broadcasts a Message to all members of a Room.
func (c *Conn) Emit(msg *Message) {
	RoomManager.Lock()
	room, ok := RoomManager.Rooms[msg.DialogID]
	RoomManager.Unlock()
	if ok == true {
		room.Emit(c, msg)
	}
}

// Upgrades an HTTP connection and creates a new Conn type.
func NewConnection(w http.ResponseWriter, r *http.Request, log *logrus.Entry,
	repo *db.Repository, userID string) *Conn {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil
	}
	c := &Conn{
		Socket: socket,
		ID:     userID,
		Send:   make(chan *Message, 10),
		Rooms:  make(map[string]string),
		log:    log,
		db:     repo,
	}
	ConnManager.Lock()
	ConnManager.Conns[c.ID] = c
	ConnManager.Unlock()
	return c
}

// Calls NewConnection, starts the returned Conn's writer, joins the root room, and finally starts the Conn's reader.
func SocketHandler(ctx echo.Context, log *logrus.Entry, repo *db.Repository, request *dto.CreateDialogRequest) error {
	c := NewConnection(ctx.Response(), ctx.Request(), log, repo, request.UserID)
	if c != nil {
		go c.writePump()
		// Запускаем для отладки рут
		c.Join("root")
		go c.readPump()
		c.log.Infof("new user: %s", c.ID)
	}
	return nil
}
