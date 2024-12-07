package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hello-world-go/internal/config"
)

func TestServerBuilder(t *testing.T) {
	appConfig := &config.AppConfig{Port: 9090}
	builder := NewServerBuilder(appConfig).
		WithRoute("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Hello, World!")); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		})

	server := builder.Build()
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	server.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("Wrong response body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestServerShutdown(t *testing.T) {
	appConfig := &config.AppConfig{Port: 9090}
	server := NewServerBuilder(appConfig).Build()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := server.Start(); err != http.ErrServerClosed {
			t.Errorf("Expected ErrServerClosed, got %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("Error during shutdown: %v", err)
	}
}
