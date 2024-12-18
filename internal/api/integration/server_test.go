package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sabasm/go-server/internal/api/handlers/health"
	"github.com/sabasm/go-server/internal/api/handlers/root"
	"github.com/sabasm/go-server/internal/server"
)

type APIResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func TestIntegrationServer(t *testing.T) {
	port, err := getAvailablePort()
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}

	cfg := &server.Config{
		Host:     "localhost",
		Port:     port,
		BasePath: "/",
	}
	srv := server.NewBuilder(cfg).
		WithRoute("/health", health.New().ServeHTTP).
		WithRoute("/", root.MustNew().ServeHTTP).
		Build()

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			t.Errorf("Server shutdown failed: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	testCases := []struct {
		name         string
		method       string
		endpoint     string
		body         interface{}
		expectedCode int
		expectedResp APIResponse
	}{
		{
			name:         "GET /health - success",
			method:       http.MethodGet,
			endpoint:     "/health",
			expectedCode: http.StatusOK,
			expectedResp: APIResponse{Status: "OK"},
		},
		{
			name:         "GET / - success",
			method:       http.MethodGet,
			endpoint:     "/",
			expectedCode: http.StatusOK,
			expectedResp: APIResponse{Message: "Welcome to Go Server"},
		},
		{
			name:         "POST /nonexistent - 404",
			method:       http.MethodPost,
			endpoint:     "/nonexistent",
			expectedCode: http.StatusNotFound,
		},
	}

	client := &http.Client{Timeout: 5 * time.Second}
	runTests(t, client, testCases, port)

	select {
	case err := <-serverErr:
		t.Fatalf("Server failed: %v", err)
	default:
	}
}

func runTests(t *testing.T, client *http.Client, testCases []struct {
	name         string
	method       string
	endpoint     string
	body         interface{}
	expectedCode int
	expectedResp APIResponse
}, port int) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			if tc.body != nil {
				body, _ = json.Marshal(tc.body)
			}

			req, err := http.NewRequest(tc.method, fmt.Sprintf("http://localhost:%d%s", port, tc.endpoint), bytes.NewReader(body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected status %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			if tc.expectedResp != (APIResponse{}) {
				var actualResp APIResponse
				if err := json.NewDecoder(resp.Body).Decode(&actualResp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if actualResp != tc.expectedResp {
					t.Errorf("Expected response %v, got %v", tc.expectedResp, actualResp)
				}
			}
		})
	}
}
