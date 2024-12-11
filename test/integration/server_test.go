package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sabasm/go-server/internal/api/handlers/health"
	"github.com/sabasm/go-server/internal/api/handlers/root"
	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/server"
)

func getAvailablePort() (int, error) {
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

func TestIntegrationServerStartup(t *testing.T) {
	port, err := getAvailablePort()
	if err != nil {
		t.Fatalf("Failed to get available port: %v", err)
	}

	configLoader := config.NewConfigLoader()
	appConfig, err := configLoader.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	appConfig.Port = port

	srv := server.NewServerBuilder(appConfig).
		WithRoute("/health", health.New().ServeHTTP).
		WithRoute("/", root.New().ServeHTTP).
		Build()

	var wg sync.WaitGroup
	wg.Add(1)
	errChan := make(chan error, 1)

	go func() {
		defer wg.Done()
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("health endpoint", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}
	})

	t.Run("root endpoint", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Server shutdown error: %v", err)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		t.Errorf("Server error: %v", err)
	default:
	}
}

func TestIntegrationHandlerResponses(t *testing.T) {
	port, err := getAvailablePort()
	if err != nil {
		t.Fatalf("Failed to get available port: %v", err)
	}

	configLoader := config.NewConfigLoader()
	appConfig, err := configLoader.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	appConfig.Port = port

	srv := server.NewServerBuilder(appConfig).
		WithRoute("/health", health.New().ServeHTTP).
		WithRoute("/", root.New().ServeHTTP).
		Build()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	testCases := []struct {
		name         string
		endpoint     string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "health check returns OK",
			endpoint:     "/health",
			expectedCode: http.StatusOK,
			expectedBody: "OK",
		},
		{
			name:         "root returns service status",
			endpoint:     "/",
			expectedCode: http.StatusOK,
			expectedBody: "Service running",
		},
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Get(fmt.Sprintf("http://localhost:%d%s", port, tc.endpoint))
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, resp.StatusCode)
			}
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Server shutdown error: %v", err)
	}

	wg.Wait()
}
