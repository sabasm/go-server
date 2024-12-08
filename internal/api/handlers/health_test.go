package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name         string
		writeError   bool
		expectedCode int
	}{
		{
			name:         "successful response",
			writeError:   false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "write error case",
			writeError:   true,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/health", nil)
			rr := httptest.NewRecorder()

			if tt.writeError {
				handler := &errorWriter{ResponseWriter: rr}
				HealthHandler(handler, req)
			} else {
				HealthHandler(rr, req)
			}

			if rr.Code != tt.expectedCode {
				t.Errorf("HealthHandler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedCode)
			}
		})
	}
}

type errorWriter struct {
	http.ResponseWriter
}

func (w *errorWriter) Write([]byte) (int, error) {
	return 0, http.ErrHandlerTimeout
}
