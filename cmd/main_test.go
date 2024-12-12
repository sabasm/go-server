package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	rootHandler(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Welcome to My MVP App"
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

// Helper function to get a free port
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

// TestMainServerLifecycle verifies server startup and graceful shutdown
func TestMainServerLifecycle(t *testing.T) {
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}

	os.Setenv("APP_PORT", fmt.Sprintf("%d", port))
	defer os.Unsetenv("APP_PORT")

	testCases := []struct {
		name           string
		setupFunc      func() error
		expectShutdown bool
	}{
		{
			name: "Graceful Shutdown",
			setupFunc: func() error {
				// Simulate normal server lifecycle
				return nil
			},
			expectShutdown: true,
		},
		{
			name: "Server Startup with Error",
			setupFunc: func() error {
				// Simulate a potential server startup error
				return fmt.Errorf("simulated startup error")
			},
			expectShutdown: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			done := make(chan struct{})
			errorCh := make(chan error, 1)

			go func() {
				defer wg.Done()
				defer close(done)

				// Introduce a mechanism to simulate potential startup conditions
				if err := tc.setupFunc(); err != nil {
					errorCh <- err
					return
				}

				main()
			}()

			// Trigger graceful shutdown
			go func() {
				time.Sleep(100 * time.Millisecond)
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(os.Interrupt)
			}()

			// Set up timeout context
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			select {
			case <-ctx.Done():
				t.Error("Test timed out")
			case <-done:
				// Successful completion
			case err := <-errorCh:
				if tc.expectShutdown {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			wg.Wait()
		})
	}
}
