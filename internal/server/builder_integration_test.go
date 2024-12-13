package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerIntegrationHTTPHandling(t *testing.T) {
	cfg := &Config{Port: 8080, Host: "localhost", BasePath: "/"}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	server := NewBuilder(cfg).
		WithRoute("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}).
		Build()

	if srv, ok := server.(*Server); ok {
		srv.router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", rec.Code)
		}
	} else {
		t.Errorf("server is not of type *Server")
	}
}
