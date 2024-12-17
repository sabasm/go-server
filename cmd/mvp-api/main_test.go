package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type mockResponseWriter struct {
	*httptest.ResponseRecorder
	writeErr error
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return m.ResponseRecorder.Write(data)
}

type mockServer struct {
	StartFunc    func() error
	ShutdownFunc func(ctx context.Context) error
}

func (m *mockServer) Start() error {
	if m.StartFunc != nil {
		return m.StartFunc()
	}
	return nil
}

func (m *mockServer) Shutdown(ctx context.Context) error {
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc(ctx)
	}
	return nil
}

func setupEnv(envVars map[string]string) {
	for k, v := range envVars {
		os.Setenv(k, v)
	}
}

func cleanupEnv(envVars map[string]string) {
	for k := range envVars {
		os.Unsetenv(k)
	}
}

func runMain(done chan struct{}) {
	go func() {
		main()
		close(done)
	}()
}

func sendRequests(t *testing.T, tests []struct {
	name           string
	path           string
	expectedCode   int
	expectedBody   map[string]string
	expectedMethod string
}) {
	client := &http.Client{Timeout: 5 * time.Second}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.expectedMethod, "http://localhost:3003"+tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			if tt.expectedBody != nil && resp.StatusCode == http.StatusOK {
				var got map[string]string
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				for k, v := range tt.expectedBody {
					if got[k] != v {
						t.Errorf("Expected %s to be %s, got %s", k, v, got[k])
					}
				}
			}
		})
	}
}

func shutdownServer(t *testing.T, proc *os.Process, done chan struct{}) {
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("Failed to send SIGTERM: %v", err)
	}

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for server shutdown")
	}
}

func TestMainFunction(t *testing.T) {
	envVars := map[string]string{
		"APP_PORT": "3003",
		"APP_HOST": "localhost",
		"DEBUG":    "true",
	}
	setupEnv(envVars)
	defer cleanupEnv(envVars)

	done := make(chan struct{})
	runMain(done)

	time.Sleep(100 * time.Millisecond)

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
		{
			name:           "Invalid Method",
			path:           "/health",
			expectedCode:   http.StatusMethodNotAllowed,
			expectedBody:   nil,
			expectedMethod: http.MethodPost,
		},
	}

	sendRequests(t, tests)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	shutdownServer(t, proc, done)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for server shutdown")
	}
}

func TestWrapHandler_FlushError(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := wrapHandler(handler, logger)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := &mockResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		writeErr:         fmt.Errorf("write error"),
	}

	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleServerShutdown_Error(t *testing.T) {
	core, obs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	mockSrv := &mockServer{
		ShutdownFunc: func(ctx context.Context) error {
			return fmt.Errorf("shutdown error")
		},
	}

	handleServerShutdown(mockSrv, logger)

	logs := obs.All()
	if len(logs) == 0 {
		t.Error("expected log entry for shutdown error")
	}

	found := false
	for _, entry := range logs {
		if entry.Message == "Forced shutdown required" {
			found = true
			if entry.Level != zap.ErrorLevel {
				t.Errorf("expected log level ERROR, got %v", entry.Level)
			}
		}
	}

	if !found {
		t.Error("expected 'Forced shutdown required' log entry")
	}
}
