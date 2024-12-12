// pkg/server/builder_test.go
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestWithTimeout(t *testing.T) {
	cfg := &Config{
		Host: "localhost",
		Port: 8080,
	}
	builder := NewBuilder(cfg)

	readTimeout := 1 * time.Second
	writeTimeout := 2 * time.Second
	idleTimeout := 3 * time.Second

	server := builder.WithTimeout(readTimeout, writeTimeout, idleTimeout).Build()

	// Type assert to access underlying server details
	if typedServer, ok := server.(*Server); ok {
		if typedServer.srv.ReadTimeout != readTimeout {
			t.Errorf("Expected ReadTimeout to be %v, got %v", readTimeout, typedServer.srv.ReadTimeout)
		}

		if typedServer.srv.WriteTimeout != writeTimeout {
			t.Errorf("Expected WriteTimeout to be %v, got %v", writeTimeout, typedServer.srv.WriteTimeout)
		}

		if typedServer.srv.IdleTimeout != idleTimeout {
			t.Errorf("Expected IdleTimeout to be %v, got %v", idleTimeout, typedServer.srv.IdleTimeout)
		}
	} else {
		t.Fatal("Could not type assert to *Server")
	}
}

func TestWithLogger(t *testing.T) {
	cfg := &Config{
		Host: "localhost",
		Port: 8080,
	}
	logger, _ := zap.NewProduction()

	builder := NewBuilder(cfg)
	server := builder.WithLogger(logger).Build()

	// Type assert to access logger
	if typedServer, ok := server.(*Server); ok {
		if typedServer.logger != logger {
			t.Error("Logger was not properly set")
		}
	} else {
		t.Fatal("Could not type assert to *Server")
	}
}

func TestWithMiddleware(t *testing.T) {
	cfg := &Config{
		Host: "localhost",
		Port: 8080,
	}
	middlewareCalled := false
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	builder := NewBuilder(cfg)
	server := builder.WithMiddleware(testMiddleware).Build()

	// Type assert to access router
	if typedServer, ok := server.(*Server); ok {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		typedServer.router.Handle("/test", handler)

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		typedServer.router.ServeHTTP(w, req)

		if !middlewareCalled {
			t.Error("Middleware was not called")
		}
	} else {
		t.Fatal("Could not type assert to *Server")
	}
}
