package chat

import (
	"sync"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

// The Dialog type represents a communication channel.
type Dialog struct {
	sync.Mutex
	Name      string
	Members   map[string]string
	stopchan  chan bool
	joinchan  chan *Conn
	leavechan chan *Conn
	Send      chan *DialogMessage
}

// Stores all Dialog types by their name.
var DialogManager = struct {
	sync.Mutex
	Rooms map[string]*Dialog
}{
	Rooms: make(map[string]*Dialog, 0),
}

// Message protocol used only with a room's Send channel.
type DialogMessage struct {
	Sender *Conn
	Data   *dto.Message
}

// Constructs and returns a new Message type.
func ConstructMessage(room, event, src, dst string, body string) *dto.Message {
	return &dto.Message{
		DialogID: room,
		Event:    event,
		AuthorID: src,
		DestinID: dst,
		Body:     body,
	}
}

// Starts the Dialog.
func (r *Dialog) Start() {
	for {
		select {
		case c := <-r.joinchan:
			var membersString string
			r.Lock()
			for id := range r.Members {
				r.Unlock()
				membersString += id + ", "
				r.Lock()
			}
			r.Members[c.ID] = c.ID
			r.Unlock()

			c.Send <- *ConstructMessage(r.Name, constants.JoinedChat, c.ID, constants.Empty, membersString)

		case c := <-r.leavechan:
			r.Lock()
			id, ok := r.Members[c.ID]
			r.Unlock()
			if ok == false {
				break
			}
			r.Lock()
			delete(r.Members, id)
			r.Unlock()
			c.Send <- *ConstructMessage(r.Name, constants.LeftChat, id, constants.Empty, c.ID)

		case rmsg := <-r.Send:
			r.Lock()

			for id := range r.Members {
				r.Unlock()

				ConnManager.Lock()
				c, ok := ConnManager.Conns[id]
				ConnManager.Unlock()
				if !ok {
					r.Lock()
					continue
				}

				select {
				case c.Send <- *rmsg.Data:
				default:
					r.Lock()
					delete(r.Members, id)
					r.Unlock()
					close(c.Send)
				}
				r.Lock()
			}

			r.Unlock()

		case <-r.stopchan:
			DialogManager.Lock()
			delete(DialogManager.Rooms, r.Name)
			DialogManager.Unlock()
			return
		}
	}
}

// Stops the Dialog.
func (r *Dialog) Stop() {
	r.stopchan <- true
}

// Adds a Conn to the Dialog.
func (r *Dialog) Join(c *Conn) {
	r.joinchan <- c
}

// Removes a Conn from the Dialog.
func (r *Dialog) Leave(c *Conn) {
	r.leavechan <- c
}

// Broadcasts data to all members of the Dialog.
func (r *Dialog) Emit(c *Conn, msg *dto.Message) {
	r.Send <- &DialogMessage{c, msg}
}

// Creates a new Dialog type and starts it.
func NewRoom(name string) *Dialog {
	r := &Dialog{
		Name:      name,
		Members:   make(map[string]string),
		stopchan:  make(chan bool),
		joinchan:  make(chan *Conn),
		leavechan: make(chan *Conn),
		Send:      make(chan *DialogMessage),
	}
	DialogManager.Lock()
	DialogManager.Rooms[name] = r
	DialogManager.Unlock()
	go r.Start()
	return r
}
