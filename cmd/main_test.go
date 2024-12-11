package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMainLifecycle(t *testing.T) {
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
			serverReady := make(chan struct{})
			serverError := make(chan error, 1)

			go func() {
				main()
				close(serverReady)
			}()

			go func() {
				time.Sleep(100 * time.Millisecond)
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(os.Interrupt)
			}()

			select {
			case <-serverReady:
			case <-time.After(2 * time.Second):
				t.Error("Server failed to start within timeout")
			case err := <-serverError:
				if !tt.expectedError {
					t.Errorf("Unexpected server error: %v", err)
				}
			}
		})
	}
}

func TestMainErrorHandling(t *testing.T) {
	os.Setenv("APP_PORT", "invalid")
	defer os.Unsetenv("APP_PORT")

	errChan := make(chan error, 1)
	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic: %v", r)
			}
			close(done)
		}()
		main()
	}()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected error for invalid port")
		}
	case <-done:
	case <-time.After(500 * time.Millisecond):
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)
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

func TestMainConfigLoading(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "8081")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
	}()

	done := make(chan struct{})
	go func() {
		main()
		close(done)
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGINT)
	}()

	select {
	case <-time.After(time.Second):
		t.Error("Test timed out")
	case <-done:
	}
}
