package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	srv := builder.Build()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK || rr.Body.String() != "Hello, World!" {
		t.Errorf("Expected 200 OK with 'Hello, World!', got %d %s", rr.Code, rr.Body.String())
	}
}



