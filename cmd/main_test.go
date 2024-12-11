package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/sabasm/go-server/internal/config"
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
		expectedError bool
	}{
		{
			name:          "graceful_shutdown",
			expectedError: false,
		},
		{
			name:          "server_error",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverError := make(chan error, 1)
			go func() {
				main()
				close(serverError)
			}()

			time.Sleep(100 * time.Millisecond)
			p, _ := os.FindProcess(os.Getpid())
			_ = p.Signal(os.Interrupt)

			select {
			case <-serverError:
			case <-time.After(2 * time.Second):
				t.Error("Server failed to start within timeout")
			}
		})
	}
}

func TestMainConfigParsing(t *testing.T) {
	os.Setenv("APP_PORT", "invalid")
	defer os.Unsetenv("APP_PORT")

	configLoader := config.NewConfigLoader()
	appConfig, err := configLoader.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	if appConfig.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", appConfig.Port)
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
	case <-done:
	case <-time.After(time.Second):
		t.Error("Test timed out")
	}
}
