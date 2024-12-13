package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestRecoveryMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Trigger-Panic") == "true" {
			panic("test panic")
		}
		w.WriteHeader(http.StatusOK)
	})

	recoveryHandler := RecoveryMiddleware(logger)(handler)

	t.Run("normal operation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		recoveryHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %v", w.Code)
		}
	})

	t.Run("panic recovery", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Trigger-Panic", "true")
		w := httptest.NewRecorder()

		recoveryHandler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %v", w.Code)
		}
	})
}
