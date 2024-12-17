package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sabasm/go-server/cmd/auth-server/service"
	"github.com/sabasm/go-server/internal/logger"
)

func TestAuthServer(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedCode   int
		expectedBody   map[string]interface{}
		expectedMethod string
	}{
		{
			name:           "M2M Token Generation",
			path:           "/auth/m2m/token",
			expectedCode:   http.StatusOK,
			expectedBody:   map[string]interface{}{"token": "test-token"},
			expectedMethod: http.MethodPost,
		},
		{
			name:           "Auth0 Token Validation",
			path:           "/auth/auth0/validate",
			expectedCode:   http.StatusOK,
			expectedBody:   map[string]interface{}{"valid": true},
			expectedMethod: http.MethodPost,
		},
	}

	os.Setenv("M2M_SECRET_KEY", "test-token")
	defer os.Unsetenv("M2M_SECRET_KEY")

	logger, _ := logger.NewLogger()
	authService := service.NewAuthService(logger)
	authService.ConfigureRoutes("test-token")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.expectedMethod, tt.path, nil)
			rr := httptest.NewRecorder()

			handler := authService.GetHandler()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			var response map[string]interface{}
			if err := decodeJSONResponse(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			for k, v := range tt.expectedBody {
				if response[k] != v {
					t.Errorf("handler returned unexpected body: got %v want %v",
						response[k], v)
				}
			}
		})
	}
}

func decodeJSONResponse(data []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewReader(data)).Decode(v)
}
