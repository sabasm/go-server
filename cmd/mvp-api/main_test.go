// cmd/mvp-api/main_test.go
package main

import (
	"context"
	"encoding/json"
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
	// Environment setup
	envVars := map[string]string{
		"APP_PORT": "3003",
		"APP_HOST": "localhost",
		"DEBUG":    "true",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	// Start server
	go func() {
		main()
	}()

	// Wait for server readiness
	time.Sleep(100 * time.Millisecond)

	// Test cases
	tests := []struct {
		name           string
		path           string
		expectedCode   int
		expectedBody   map[string]string
		expectedMethod string
	}{
		{
			name:           "Health Check Endpoint",
			path:           "/health",
			expectedCode:   http.StatusOK,
			expectedBody:   map[string]string{"status": "OK"},
			expectedMethod: http.MethodGet,
		},
		{
			name:           "Root Endpoint",
			path:           "/",
			expectedCode:   http.StatusOK,
			expectedBody:   map[string]string{"message": "Welcome to Go Server"},
			expectedMethod: http.MethodGet,
		},
		{
			name:           "Not Found Endpoint",
			path:           "/notfound",
			expectedCode:   http.StatusNotFound,
			expectedMethod: http.MethodGet,
		},
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.expectedMethod,
				"http://localhost:3003"+tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d",
					tt.expectedCode, resp.StatusCode)
			}

			if tt.expectedBody != nil {
				var got map[string]string
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				for k, v := range tt.expectedBody {
					if got[k] != v {
						t.Errorf("Expected %s to be %s, got %s",
							k, v, got[k])
					}
				}
			}
		})
	}

	// Graceful shutdown
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("Failed to send SIGTERM: %v", err)
	}
}

func TestHandleServerShutdown(t *testing.T) {
	cfg := &server.Config{
		Host:     "localhost",
		Port:     3004,
		BasePath: "/",
		Options: server.Options{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}

	srv := server.NewBuilder(cfg).
		WithRoute("/health", health.New().ServeHTTP).
		WithRoute("/", root.New().ServeHTTP).
		Build()

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	t.Run("Validates Graceful Shutdown", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown failed: %v", err)
		}

		select {
		case err := <-serverErr:
			t.Errorf("Unexpected server error: %v", err)
		default:
		}
	})
}
