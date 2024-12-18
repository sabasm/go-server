package root

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorWriter struct {
	http.ResponseWriter
}

func (ew *errorWriter) Write([]byte) (int, error) {
	return 0, http.ErrHandlerTimeout
}

func TestRootHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		wantStatus   int
		wantResponse RootResponse
		forceError   bool
	}{
		{
			name:         "GET Method",
			method:       http.MethodGet,
			wantStatus:   http.StatusOK,
			wantResponse: RootResponse{Message: welcomeMsg},
		},
		{
			name:         "POST Method",
			method:       http.MethodPost,
			wantStatus:   http.StatusOK,
			wantResponse: RootResponse{Message: welcomeMsg},
		},
		{
			name:         "PUT Method",
			method:       http.MethodPut,
			wantStatus:   http.StatusOK,
			wantResponse: RootResponse{Message: welcomeMsg},
		},
		{
			name:         "DELETE Method",
			method:       http.MethodDelete,
			wantStatus:   http.StatusOK,
			wantResponse: RootResponse{Message: welcomeMsg},
		},
		{
			name:         "PATCH Method",
			method:       http.MethodPatch,
			wantStatus:   http.StatusOK,
			wantResponse: RootResponse{Message: welcomeMsg},
		},
		{
			name:         "Invalid Method",
			method:       http.MethodTrace,
			wantStatus:   http.StatusMethodNotAllowed,
			wantResponse: RootResponse{Error: methodNotAllowedMsg},
		},
		{
			name:       "Write Error - OK Response",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			forceError: true,
		},
		{
			name:       "Write Error - Method Not Allowed",
			method:     http.MethodTrace,
			wantStatus: http.StatusMethodNotAllowed,
			forceError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := MustNew()
			req := httptest.NewRequest(tt.method, "/", nil)
			rec := httptest.NewRecorder()

			var w http.ResponseWriter = rec
			if tt.forceError {
				w = &errorWriter{rec}
			}

			handler.ServeHTTP(w, req)

			if got := rec.Code; got != tt.wantStatus {
				t.Errorf("status code = %d, want %d", got, tt.wantStatus)
			}

			if !tt.forceError {
				var response RootResponse
				if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response != tt.wantResponse {
					t.Errorf("response = %+v, want %+v", response, tt.wantResponse)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	handler, err := New()
	if err != nil {
		t.Fatalf("New() returned unexpected error: %v", err)
	}
	if handler.logger == nil {
		t.Error("New() returned handler with nil logger")
	}
}

func TestMustNew(t *testing.T) {
	handler := MustNew()
	if handler == nil {
		t.Fatal("MustNew() returned nil")
	}
	if handler.logger == nil {
		t.Error("MustNew() returned handler with nil logger")
	}
}
