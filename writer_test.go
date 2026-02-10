package binproto_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tetsuo/binproto"
)

func TestSend(t *testing.T) {
	var buf bytes.Buffer
	w := binproto.NewWriter(bufio.NewWriter(&buf))
	err := w.WriteMessage(newMessage(42, 3, 2))
	if s := buf.String(); s != "\x04\xa3\x05ab" || err != nil {
		t.Fatalf("s=%q; err=%s", s, err)
	}
}

func TestSendBatch(t *testing.T) {
	var buf bytes.Buffer
	w := binproto.NewWriter(bufio.NewWriter(&buf))
	msg := newMessage(42, 3, 2)
	err := w.WriteMessage(msg, msg)
	if s := buf.String(); s != "\x04\xa3\x05ab\x04\xa3\x05ab" || err != nil {
		t.Fatalf("s=%q; err=%s", s, err)
	}
}

func TestWriteMessageNilData(t *testing.T) {
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	w := binproto.NewWriter(bw)

	msg := binproto.NewMessage(100, 5, nil)
	err := w.WriteMessage(msg)
	assert.NoError(t, err)

	r := binproto.NewReaderSize(bytes.NewReader(buf.Bytes()), 64)
	m := &binproto.Message{}
	err = r.ReadMessage(m)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), m.ID)
	assert.Equal(t, uint8(5), m.Type)
	assert.Equal(t, 0, len(m.Data))
}

func TestConnClose(t *testing.T) {
	type mockConn struct {
		*bytes.Buffer
		closed bool
	}

	buf := &mockConn{Buffer: &bytes.Buffer{}}

	data := []byte{0x04, 0xa3, 0x05, 'a', 'b'}
	buf.Write(data)

	conn := binproto.NewConn(struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		Reader: buf,
		Writer: buf,
		Closer: io.NopCloser(buf),
	})

	err := conn.Close()
	assert.NoError(t, err)
}

// Test Send error path when WriteMessage fails
type failWriter struct {
	*bytes.Buffer
}

func (f *failWriter) Close() error {
	return nil
}

func (f *failWriter) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (f *failWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestSendWriteError(t *testing.T) {
	fw := &failWriter{Buffer: &bytes.Buffer{}}
	conn := binproto.NewConn(fw)

	msg := binproto.NewMessage(42, 3, []byte("test"))
	_, err := conn.Send(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write failed")
}

func TestDialError(t *testing.T) {
	_, err := binproto.Dial("tcp", "invalid:99999999")
	assert.Error(t, err)
}
