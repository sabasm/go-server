package handlers

import (
	"bytes"
	"net/http"
)

type BufferedResponseWriter struct {
	http.ResponseWriter
	buffer     *bytes.Buffer
	statusCode int
}

func NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		ResponseWriter: w,
		buffer:         new(bytes.Buffer),
		statusCode:     http.StatusOK,
	}
}

func (b *BufferedResponseWriter) WriteHeader(code int) {
	b.statusCode = code
}

func (b *BufferedResponseWriter) Write(data []byte) (int, error) {
	return b.buffer.Write(data)
}

func (b *BufferedResponseWriter) Flush() error {
	if b.statusCode != http.StatusOK {
		b.ResponseWriter.WriteHeader(b.statusCode)
	}
	_, err := b.ResponseWriter.Write(b.buffer.Bytes())
	return err
}
