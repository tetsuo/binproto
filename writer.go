package binproto

import (
	"bufio"
	"encoding/binary"
)

// A Writer implements convenience methods for writing
// requests or responses to a binary protocol network connection.
type Writer struct {
	wd *bufio.Writer
}

// NewWriter returns a new Writer writing to w.
func NewWriter(wd *bufio.Writer) *Writer {
	return &Writer{wd: wd}
}

// WriteMessage writes a variable number of messages to w.
func (w *Writer) WriteMessage(messages ...*Message) error {
	var vb [10]byte // max 10 bytes per varint

	for _, m := range messages {
		header := uint64(m.ID<<4) | uint64(m.Channel)
		bodyLen := uint64(encodingLength(header) + len(m.Data))

		// write body length varint
		n := binary.PutUvarint(vb[:], bodyLen)
		if _, err := w.wd.Write(vb[:n]); err != nil {
			return err
		}

		// write header varint
		n = binary.PutUvarint(vb[:], header)
		if _, err := w.wd.Write(vb[:n]); err != nil {
			return err
		}

		// write data payload
		if _, err := w.wd.Write(m.Data); err != nil {
			return err
		}
	}

	return w.wd.Flush()
}
