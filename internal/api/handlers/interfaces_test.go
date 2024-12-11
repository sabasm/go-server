package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandler struct {
	serveHTTPCalled bool
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.serveHTTPCalled = true
}

func TestHandlerInterface(t *testing.T) {
	mock := &mockHandler{}
	var _ Handler = mock

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	mock.ServeHTTP(w, r)

	if !mock.serveHTTPCalled {
		t.Error("ServeHTTP was not called")
	}
}
