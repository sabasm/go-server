package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testPayload struct {
	Message string `json:"message"`
}

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		payload    interface{}
		wantStatus int
		wantBody   string
	}{
		{
			name:       "struct payload",
			payload:    testPayload{Message: "test"},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"test"}`,
		},
		{
			name:       "map payload",
			payload:    map[string]string{"key": "value"},
			wantStatus: http.StatusCreated,
			wantBody:   `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := WriteJSONResponse(w, tt.wantStatus, tt.payload)
			if err != nil {
				t.Fatalf("WriteJSONResponse() error = %v", err)
			}

			if status := w.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		})
	}
}
