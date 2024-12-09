package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMainServerLifecycle(t *testing.T) {
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}

	os.Setenv("APP_PORT", fmt.Sprintf("%d", port))
	defer os.Unsetenv("APP_PORT")

	tests := []struct {
		name          string
		injectError   bool
		expectedError bool
	}{
		{
			name:          "graceful_shutdown",
			injectError:   false,
			expectedError: false,
		},
		{
			name:          "server_error",
			injectError:   true,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				if tt.injectError {
					mockServer := &http.Server{
						Addr:              ":invalid",
						ReadHeaderTimeout: 3 * time.Second,
					}
					_ = mockServer.ListenAndServe()
				}
			}()

			go func() {
				time.Sleep(100 * time.Millisecond)
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(os.Interrupt)
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			done := make(chan struct{})
			go func() {
				main()
				close(done)
			}()

			select {
			case <-ctx.Done():
				t.Error("Test timed out")
			case <-done:
			}

			wg.Wait()
		})
	}
}

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
