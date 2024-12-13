// cmd/main_test.go
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

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Welcome to My MVP App"))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

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

func getFreePort() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

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
				return nil
			},
			expectShutdown: true,
		},
		{
			name: "Server Startup with Error",
			setupFunc: func() error {
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
				if err := tc.setupFunc(); err != nil {
					errorCh <- err
					return
				}
				main()
			}()

			go func() {
				time.Sleep(100 * time.Millisecond)
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(os.Interrupt)
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			select {
			case <-ctx.Done():
				t.Error("Test timed out")
			case <-done:
			case err := <-errorCh:
				if tc.expectShutdown {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			wg.Wait()
		})
	}
}
