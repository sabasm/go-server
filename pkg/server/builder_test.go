package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerBuilder(t *testing.T) {
	cfg := &Config{Port: 8080, Host: "localhost"}
	tests := []struct {
		name     string
		build    func(ServerBuilder) ServerBuilder
		validate func(*testing.T, *Server)
	}{
		{
			name: "basic_server",
			build: func(b ServerBuilder) ServerBuilder {
				return b
			},
			validate: func(t *testing.T, s *Server) {
				if s.srv.Addr != "localhost:8080" {
					t.Errorf("expected addr localhost:8080, got %s", s.srv.Addr)
				}
			},
		},
		{
			name: "with_timeouts",
			build: func(b ServerBuilder) ServerBuilder {
				return b.WithTimeout(1*time.Second, 2*time.Second, 3*time.Second)
			},
			validate: func(t *testing.T, s *Server) {
				if s.srv.ReadTimeout != time.Second {
					t.Errorf("expected read timeout 1s, got %v", s.srv.ReadTimeout)
				}
			},
		},
		{
			name: "with_route",
			build: func(b ServerBuilder) ServerBuilder {
				return b.WithRoute("/test", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
			},
			validate: func(t *testing.T, s *Server) {
				req := httptest.NewRequest("GET", "/test", nil)
				rr := httptest.NewRecorder()
				s.router.ServeHTTP(rr, req)
				if rr.Code != http.StatusOK {
					t.Errorf("expected status OK, got %v", rr.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(cfg)
			srv := tt.build(builder).Build()
			tt.validate(t, srv)
		})
	}
}
