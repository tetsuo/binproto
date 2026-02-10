package binproto

// A Message represents a single binproto message.
//
// Each message starts with an header which is a varint encoded
// unsigned 64-bit integer which consists of a channel ID (first 60-bits) and
// a message type (last 4-bits), the rest of the message is payload.
type Message struct {
	ID   uint64
	Type uint8
	Data []byte
}

// NewMessage returns a new Message.
func NewMessage(id uint64, messageType uint8, data []byte) *Message {
	return &Message{
		ID:   id,
		Type: messageType,
		Data: data,
	}
}
