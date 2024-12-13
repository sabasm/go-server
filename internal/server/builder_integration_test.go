package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerIntegration(t *testing.T) {
	cfg := &Config{Port: 8080, Host: "localhost", BasePath: "/"}

	t.Run("server lifecycle", func(t *testing.T) {
		srv := NewBuilder(cfg).
			WithRoute("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}).
			Build()

		go func() {
			if err := srv.Start(); err != nil && err != http.ErrServerClosed {
				t.Errorf("unexpected error: %v", err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			t.Errorf("shutdown error: %v", err)
		}
	})

	t.Run("http handling", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		server := NewBuilder(cfg).
			WithRoute("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}).
			Build()

		if srv, ok := server.(*Server); ok {
			srv.router.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("expected status OK; got %v", rec.Code)
			}
		}
	})
}
