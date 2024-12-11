package root

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorWriter struct {
	http.ResponseWriter
	forceError bool
}

func (w *errorWriter) Write([]byte) (int, error) {
	if w.forceError {
		return 0, errors.New("forced write error")
	}
	return len([]byte("Service running")), nil
}

func TestRootHandler_ServeHTTP(t *testing.T) {
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
			handler := New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rw := httptest.NewRecorder()
			writer := &errorWriter{ResponseWriter: rw, forceError: tt.writeError}

			handler.ServeHTTP(writer, req)

			if rw.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rw.Code, tt.expectedCode)
			}
		})
	}
}
