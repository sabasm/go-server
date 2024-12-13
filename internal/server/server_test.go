package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestServerLifecycle(t *testing.T) {
	cfg := &Config{
		Port: 0, // Random port for testing
		Host: "localhost",
	}

	srv := NewBuilder(cfg).
		WithRoute("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}).
		Build()

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("server shutdown failed: %v", err)
	}

	select {
	case err := <-errCh:
		t.Fatalf("server error: %v", err)
	default:
	}
}
