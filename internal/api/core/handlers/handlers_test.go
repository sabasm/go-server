package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHandler struct {
	called bool
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	w.WriteHeader(http.StatusOK)
}

func TestHandlerInterface(t *testing.T) {
	t.Run("handler implementation", func(t *testing.T) {
		var _ Handler = &testHandler{}

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
	})

	t.Run("error handling", func(t *testing.T) {
		h := &testHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("Trigger-Error", "true")

		h.ServeHTTP(w, r)

		if !h.called {
			t.Error("handler was not called")
		}
	})
}
