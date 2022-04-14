package chat

import (
	"encoding/json"
	"log"
	"sync"
)

// The Room type represents a communication channel.
type Room struct {
	sync.Mutex
	ID        string
	Members   map[string]string
	stopchan  chan bool
	joinchan  chan *Conn
	leavechan chan *Conn
	Send      chan *RoomMessage
}

// Stores all Room types by their name.
var RoomManager = struct {
	sync.Mutex
	Rooms map[string]*Room
}{
	Rooms: make(map[string]*Room, 0),
}

// Starts the Room.
func (r *Room) Start() {
	for {
		select {
		case c := <-r.joinchan:
			members := make([]string, 0)
			r.Lock()
			for id := range r.Members {
				r.Unlock()
				members = append(members, id)
				r.Lock()
			}
			r.Members[c.ID] = c.ID
			r.Unlock()
			payload, err := json.Marshal(members)
			if err != nil {
				log.Println(err)
				break
			}
			c.Send <- ConstructMessage(r.ID, "join", "", c.ID, payload)
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
			c.log.Infof("ConstructMessage, left chat")
			c.Send <- ConstructMessage(r.ID, "leave", "", id, []byte(c.ID))
		case rmsg := <-r.Send:
			rmsg.Sender.log.Infof("Sendner %s, event: %s, payload: %s",
				rmsg.Sender.ID, rmsg.Data.Event, rmsg.Data.Payload)
			r.Lock()
			for id := range r.Members {
				r.Unlock()
				ConnManager.Lock()
				c, ok := ConnManager.Conns[id]
				ConnManager.Unlock()
				if !ok || c == rmsg.Sender {
					r.Lock()
					continue
				}
				select {
				case c.Send <- rmsg.Data:
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
			RoomManager.Lock()
			delete(RoomManager.Rooms, r.ID)
			RoomManager.Unlock()
			return
		}
	}
}

// Stops the Room.
func (r *Room) Stop() {
	r.stopchan <- true
}

// Adds a Conn to the Room.
func (r *Room) Join(c *Conn) {
	r.joinchan <- c
}

// Removes a Conn from the Room.
func (r *Room) Leave(c *Conn) {
	r.leavechan <- c
}

// Broadcasts data to all members of the Room.
func (r *Room) Emit(c *Conn, msg *Message) {
	r.Send <- &RoomMessage{c, msg}
}

// Creates a new Room type and starts it.
func NewRoom(id string) *Room {
	r := &Room{
		ID:        id,
		Members:   make(map[string]string),
		stopchan:  make(chan bool),
		joinchan:  make(chan *Conn),
		leavechan: make(chan *Conn),
		Send:      make(chan *RoomMessage),
	}
	RoomManager.Lock()
	RoomManager.Rooms[id] = r
	RoomManager.Unlock()
	go r.Start()
	return r
}
