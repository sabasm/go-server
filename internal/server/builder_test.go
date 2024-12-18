package server

import (
	"net/http"
	"testing"

	"go.uber.org/zap"
)

func TestWithMiddleware(t *testing.T) {
	cfg := &Config{Host: "localhost", Port: 8080}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "middleware")
			next.ServeHTTP(w, r)
		})
	}

	builder := NewBuilder(cfg).WithMiddleware(middleware)

	if len(builder.middleware) != 1 {
		t.Fatalf("expected 1 middleware, got %d", len(builder.middleware))
	}
}

func TestWithLogger(t *testing.T) {
	cfg := &Config{Host: "localhost", Port: 8080}
	logger, _ := zap.NewDevelopment()

	builder := NewBuilder(cfg).WithLogger(logger)

	if builder.logger == nil {
		t.Fatal("expected logger to be set")
	}
}

func TestBuildWithNilConfig(t *testing.T) {
	builder := NewBuilder(nil).Build()

	if builder != nil {
		t.Fatal("expected nil server when config is nil")
	}
}

func TestBuildWithInvalidRouter(t *testing.T) {
	builder := &serverBuilder{config: &Config{Host: "localhost", Port: 8080}}
	server := builder.Build()

	if server != nil {
		t.Fatal("expected nil server when router is uninitialized")
	}
}
