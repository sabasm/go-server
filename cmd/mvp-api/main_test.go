// cmd/mvp-api/main_test.go
package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/sabasm/go-server/internal/api/handlers/health"
	"github.com/sabasm/go-server/internal/api/handlers/root"
	"github.com/sabasm/go-server/internal/server"
)

func TestMainFunction(t *testing.T) {
	// Override environment variables
	os.Setenv("APP_PORT", "3003")
	defer os.Unsetenv("APP_PORT")

	// Run the main function in a separate goroutine
	go func() {
		main()
	}()

	// Allow server time to start
	time.Sleep(100 * time.Millisecond)

	// Test Health Check Endpoint
	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := http.Get("http://localhost:3003/health")
		if err != nil {
			t.Fatalf("Failed to reach health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test Root Endpoint
	t.Run("RootHandler", func(t *testing.T) {
		resp, err := http.Get("http://localhost:3003/")
		if err != nil {
			t.Fatalf("Failed to reach root endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Trigger SIGTERM to gracefully shut down
	proc, _ := os.FindProcess(os.Getpid())
	_ = proc.Signal(syscall.SIGTERM)

	// Allow shutdown process
	time.Sleep(200 * time.Millisecond)
}

func TestHandleServerShutdown(t *testing.T) {
	cfg := &server.Config{
		Host:     "localhost",
		Port:     3004,
		BasePath: "/",
	}

	srv := server.NewBuilder(cfg).
		WithRoute("/health", health.New().ServeHTTP).
		WithRoute("/", root.New().ServeHTTP).
		Build()

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server start error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("GracefulShutdown", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			t.Errorf("Server shutdown failed: %v", err)
		}
	})
}
