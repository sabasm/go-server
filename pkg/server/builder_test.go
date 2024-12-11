package server

import (
	"testing"
	"time"
)

func TestServerBuilder(t *testing.T) {
	cfg := &Config{Port: 8080, Host: "localhost"}
	tests := []struct {
		name     string
		build    func(ServerBuilder) ServerBuilder
		validate func(*testing.T, ServerInterface)
	}{
		{
			name: "basic_server",
			build: func(b ServerBuilder) ServerBuilder {
				return b
			},
			validate: func(t *testing.T, s ServerInterface) {
				if srv, ok := s.(*Server); !ok {
					t.Error("expected Server type")
				} else if srv.srv.Addr != "localhost:8080" {
					t.Errorf("expected addr localhost:8080, got %s", srv.srv.Addr)
				}
			},
		},
		{
			name: "with_timeouts",
			build: func(b ServerBuilder) ServerBuilder {
				return b.WithTimeout(1*time.Second, 2*time.Second, 3*time.Second)
			},
			validate: func(t *testing.T, s ServerInterface) {
				if srv, ok := s.(*Server); !ok {
					t.Error("expected Server type")
				} else if srv.srv.ReadTimeout != time.Second {
					t.Errorf("expected read timeout 1s, got %v", srv.srv.ReadTimeout)
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
