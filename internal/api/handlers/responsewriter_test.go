// pkg/api/handlers/responsewriter_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBufferedResponseWriter(t *testing.T) {
	tests := []struct {
		name        string
		writeData   []byte
		statusCode  int
		expectError bool
	}{
		{
			name:        "successful write",
			writeData:   []byte("test data"),
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "custom status code",
			writeData:   []byte("error data"),
			statusCode:  http.StatusBadRequest,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			bw := NewBufferedResponseWriter(rr)

			n, err := bw.Write(tt.writeData)
			if err != nil != tt.expectError {
				t.Errorf("Write() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if n != len(tt.writeData) {
				t.Errorf("Write() wrote %d bytes; want %d", n, len(tt.writeData))
			}

			bw.WriteHeader(tt.statusCode)
			if err := bw.Flush(); err != nil != tt.expectError {
				t.Errorf("Flush() error = %v, expectError %v", err, tt.expectError)
			}

			if rr.Code != tt.statusCode {
				t.Errorf("status code = %d, want %d", rr.Code, tt.statusCode)
			}
			if rr.Body.String() != string(tt.writeData) {
				t.Errorf("body = %q, want %q", rr.Body.String(), string(tt.writeData))
			}
		})
	}
}
