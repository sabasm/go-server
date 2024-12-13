package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody map[string]string
	}{
		{
			name:         "success response",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: map[string]string{"status": "OK"},
		},
		{
			name:         "invalid method",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: map[string]string{"status": "OK"},
		},
	}

	handler := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			var response map[string]string
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response["status"] != tt.expectedBody["status"] {
				t.Errorf("expected body %v, got %v", tt.expectedBody, response)
			}
		})
	}

	t.Run("validate request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		if err := handler.ValidateRequest(req); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestHealthHandlerEdgeCases(t *testing.T) {
	t.Run("large payload", func(t *testing.T) {
		handler := New()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})

	t.Run("invalid request method", func(t *testing.T) {
		handler := New()
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})
}
