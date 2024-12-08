// test/integration/server_test.go
package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

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
	port, portErr := getAvailablePort()
	if portErr != nil {
		t.Fatalf("Failed to get available port: %v", portErr)
	}

	configLoader := config.NewConfigLoader()
	appConfig, configErr := configLoader.LoadConfig()
	if configErr != nil {
		t.Fatalf("Failed to load config: %v", configErr)
	}
	appConfig.Port = port

	srv := server.NewServerBuilder(appConfig).
		WithRoute("/test", func(w http.ResponseWriter, r *http.Request) {
			response := []byte("test response")
			n, writeErr := w.Write(response)
			if writeErr != nil {
				t.Errorf("Failed to write response: %v", writeErr)
				return
			}
			if n != len(response) {
				t.Errorf("Failed to write complete response: wrote %d bytes of %d", n, len(response))
			}
		}).
		Build()

	var wg sync.WaitGroup
	wg.Add(1)
	errChan := make(chan error, 1)

	go func() {
		defer wg.Done()
		startErr := srv.Start()
		if startErr != nil && startErr != http.ErrServerClosed {
			errChan <- startErr
		}
	}()

	time.Sleep(100 * time.Millisecond)

	resp, httpErr := http.Get(fmt.Sprintf("http://localhost:%d/test", port))
	if httpErr != nil {
		t.Fatalf("Failed to make request: %v", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownErr := srv.Shutdown(ctx)
	if shutdownErr != nil {
		t.Errorf("Server shutdown error: %v", shutdownErr)
	}

	wg.Wait()

	select {
	case serverErr := <-errChan:
		t.Errorf("Server error: %v", serverErr)
	default:
	}
}
