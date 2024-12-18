package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	core, obs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Change the request method to POST to trigger logging
	req := httptest.NewRequest("POST", "/test", nil)
	rr := httptest.NewRecorder()

	mw := LoggingMiddleware(logger)
	mw(handler).ServeHTTP(rr, req)

	logs := obs.All()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logs))
	}

	if logs[0].Message != "request" {
		t.Errorf("expected message 'request', got '%s'", logs[0].Message)
	}

	// Optionally, verify log fields
	expectedMethod := "POST"
	expectedPath := "/test"

	if logs[0].ContextMap()["method"] != expectedMethod {
		t.Errorf("expected method '%s', got '%s'", expectedMethod, logs[0].ContextMap()["method"])
	}

	if logs[0].ContextMap()["path"] != expectedPath {
		t.Errorf("expected path '%s', got '%s'", expectedPath, logs[0].ContextMap()["path"])
	}
}
