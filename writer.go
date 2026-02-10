package binproto

import (
	"bufio"
	"encoding/binary"
)

// A Writer implements convenience methods for writing
// requests or responses to a binary protocol network connection.
type Writer struct {
	wd      *bufio.Writer
	scratch [20]byte
}

// NewWriter returns a new Writer writing to w.
func NewWriter(wd *bufio.Writer) *Writer {
	return &Writer{wd: wd}
}

// WriteMessage writes a variable number of messages to w.
func (w *Writer) WriteMessage(messages ...*Message) error {
	for _, m := range messages {
		header := m.ID<<4 | uint64(m.Type)
		headerLen := binary.PutUvarint(w.scratch[10:], header)
		bodyLen := uint64(headerLen + len(m.Data))
		lengthSize := binary.PutUvarint(w.scratch[:10], bodyLen)
		// Shift header bytes to be adjacent to length bytes
		copy(w.scratch[lengthSize:], w.scratch[10:10+headerLen])
		if _, err := w.wd.Write(w.scratch[:lengthSize+headerLen]); err != nil {
			return err
		}
		if len(m.Data) > 0 {
			if _, err := w.wd.Write(m.Data); err != nil {
				return err
			}
		}
	}
	return w.wd.Flush()
}
