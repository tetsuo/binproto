package binproto_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/tetsuo/binproto"
)

func BenchmarkWriter_64B(b *testing.B) {
	benchWriter(b, 64)
}

func BenchmarkWriter_4KB(b *testing.B) {
	benchWriter(b, 4*1024)
}

func BenchmarkWriter_1MB(b *testing.B) {
	benchWriter(b, 1024*1024)
}

func benchWriter(b *testing.B, payloadSize int) {
	payload := make([]byte, payloadSize)
	msg := binproto.NewMessage(1, 'X', payload)

	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	w := binproto.NewWriter(bw)
	b.SetBytes(int64(payloadSize))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		bw.Reset(buf)
		w.WriteMessage(msg)
		bw.Flush()
	}
}

func BenchmarkReader_64B(b *testing.B) {
	benchReader(b, 64)
}

func BenchmarkReader_4KB(b *testing.B) {
	benchReader(b, 4*1024)
}

func BenchmarkReader_1MB(b *testing.B) {
	benchReader(b, 1024*1024)
}

func benchReader(b *testing.B, payloadSize int) {
	// Encode message once
	payload := make([]byte, payloadSize)
	msg := binproto.NewMessage(1, 'X', payload)
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	w := binproto.NewWriter(bw)
	w.WriteMessage(msg)
	bw.Flush()
	encoded := buf.Bytes()

	// Buffer size must be larger than encoded message size
	bufSize := 32 * 1024
	if len(encoded) > bufSize {
		bufSize = len(encoded) + 1024 // Add some headroom
	}

	source := &repeatReader{data: encoded}
	reader := binproto.NewReaderSize(source, bufSize)

	b.SetBytes(int64(payloadSize))
	b.ReportAllocs()
	b.ResetTimer()

	m := &binproto.Message{}

	for i := 0; i < b.N; i++ {
		err := reader.ReadMessage(m)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// repeatReader repeats the same data indefinitely.
type repeatReader struct {
	data []byte
	pos  int
}

func (r *repeatReader) Read(p []byte) (n int, err error) {
	if len(r.data) == 0 {
		return 0, nil
	}

	for n < len(p) && n < len(r.data) {
		toCopy := len(p) - n
		remaining := len(r.data) - r.pos
		if toCopy > remaining {
			toCopy = remaining
		}

		copy(p[n:], r.data[r.pos:r.pos+toCopy])
		n += toCopy
		r.pos += toCopy

		if r.pos >= len(r.data) {
			r.pos = 0
		}
	}

	return n, nil
}
