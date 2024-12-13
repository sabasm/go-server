package handlers_test

import (
        "net/http"
        "net/http/httptest"
        "testing"

        "github.com/sabasm/go-server/internal/api/core/handlers"
)

type testHandler struct {
        called bool
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        h.called = true
        w.WriteHeader(http.StatusOK)
}

func TestHandlerInterface(t *testing.T) {
        var _ handlers.Handler = &testHandler{}

        h := &testHandler{}
        w := httptest.NewRecorder()
        r := httptest.NewRequest(http.MethodGet, "/test", nil)

        h.ServeHTTP(w, r)

        if !h.called {
                t.Error("handler was not called")
        }

        if w.Code != http.StatusOK {
                t.Errorf("expected status OK, got %v", w.Code)
        }
}


