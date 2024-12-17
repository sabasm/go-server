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
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: nil,
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

			if tt.expectedBody != nil {
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				for k, v := range tt.expectedBody {
					if got, exists := response[k]; !exists || got != v {
						t.Errorf("expected %s to be %s, got %s", k, v, got)
					}
				}
			}
		})
	}
}
