package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/server"
	"github.com/sabasm/go-server/pkg/api/handlers/health"
	"github.com/sabasm/go-server/pkg/api/handlers/root"
)

type testResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
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
		expectStatus string
		expectMsg    string
	}{
		{
			name:         "health check returns OK",
			endpoint:     "/health",
			expectStatus: "OK",
		},
		{
			name:      "root returns welcome message",
			endpoint:  "/",
			expectMsg: "Welcome to Go Server",
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

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK, got %v", resp.Status)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			var response testResponse
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tc.expectStatus != "" && response.Status != tc.expectStatus {
				t.Errorf("Expected status %q, got %q", tc.expectStatus, response.Status)
			}

			if tc.expectMsg != "" && response.Message != tc.expectMsg {
				t.Errorf("Expected message %q, got %q", tc.expectMsg, response.Message)
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
