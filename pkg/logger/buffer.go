package logger

import (
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type bufferWriter struct {
	buffer *ringbuf.Buffer[string]
}

func (w *bufferWriter) Write(p []byte) (n int, err error) {
	w.buffer.Push(string(p))
	return len(p), nil
}

func RegisterBufferLogger(bufSize int) *ringbuf.Buffer[string] {
	buf := ringbuf.New[string](bufSize)
	logWriters = append(logWriters, &bufferWriter{buffer: buf})
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(logWriters...)).With().Timestamp().Logger()
	return buf
}
