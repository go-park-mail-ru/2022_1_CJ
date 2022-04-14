package chat

type Message struct {
	DialogName string `json:"dialog_name" bson:"dialog_name"`
	DialogID   string `json:"dialog_id" bson:"dialog_id"`
	Event      string `json:"type" bson:"event"`
	SrcID      string `json:"src_id" bson:"src_id"`
	Payload    []byte `json:"payload" bson:"payload"`
}

// Message protocol used only with a room's Send channel.
type RoomMessage struct {
	Sender *Conn    `json:"sender"`
	Data   *Message `json:"data"`
}

// Constructs and returns a new Message type.
func ConstructMessage(room, event, dst, src string, payload []byte) *Message {
	return &Message{
		DialogName: room,
		DialogID:   dst,
		Event:      event,
		SrcID:      src,
		Payload:    payload,
	}
}
