package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestHandleServerShutdown_Error(t *testing.T) {
	s := &Server{
		srv: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: zap.NewNop(),
	}

	err := s.Shutdown(context.Background())
	if err != nil {
		t.Errorf("unexpected shutdown error: %v", err)
	}
}

func TestServerWithInvalidConfig(t *testing.T) {
	builder := NewBuilder(nil)
	if builder != nil {
		t.Fatal("expected nil server when config is nil")
	}
}
