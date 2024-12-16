package server

import (
	"net/http"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestBuilderMethods(t *testing.T) {
	validConfig := &Config{
		Host: "localhost",
		Port: 8080,
		Options: Options{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}

	logger, _ := zap.NewDevelopment()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := func(next http.Handler) http.Handler { return next }

	tests := []struct {
		name      string
		builder   *serverBuilder
		operation func(*serverBuilder) *serverBuilder
		checkFunc func(*testing.T, *serverBuilder)
		wantNil   bool
	}{
		{
			name:    "WithRoute",
			builder: NewBuilder(validConfig),
			operation: func(b *serverBuilder) *serverBuilder {
				return b.WithRoute("/test", handler)
			},
			checkFunc: func(t *testing.T, b *serverBuilder) {
				if b.router == nil {
					t.Error("router should not be nil")
				}
			},
		},
		{
			name:    "WithMiddleware",
			builder: NewBuilder(validConfig),
			operation: func(b *serverBuilder) *serverBuilder {
				return b.WithMiddleware(middleware)
			},
			checkFunc: func(t *testing.T, b *serverBuilder) {
				if len(b.middleware) != 1 {
					t.Error("middleware should be added")
				}
			},
		},
		{
			name:    "WithLogger",
			builder: NewBuilder(validConfig),
			operation: func(b *serverBuilder) *serverBuilder {
				return b.WithLogger(logger)
			},
			checkFunc: func(t *testing.T, b *serverBuilder) {
				if b.logger == nil {
					t.Error("logger should be set")
				}
			},
		},
		{
			name:    "WithTimeout",
			builder: NewBuilder(validConfig),
			operation: func(b *serverBuilder) *serverBuilder {
				return b.WithTimeout(1*time.Second, 1*time.Second, 1*time.Second)
			},
			checkFunc: func(t *testing.T, b *serverBuilder) {
				if b.timeouts.read != time.Second {
					t.Error("timeout should be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operation(tt.builder)
			if !tt.wantNil {
				if result == nil {
					t.Fatal("expected non-nil result")
				}
				tt.checkFunc(t, result)
			}
		})
	}
}

func TestBuildWithValidations(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		setup     func(*serverBuilder)
		wantError bool
	}{
		{
			name: "complete_valid_setup",
			config: &Config{
				Host: "localhost",
				Port: 8080,
				Options: Options{
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 5 * time.Second,
					IdleTimeout:  30 * time.Second,
				},
			},
			setup: func(b *serverBuilder) {
				b.WithRoute("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			},
			wantError: false,
		},
		{
			name:      "nil_config",
			config:    nil,
			setup:     func(b *serverBuilder) {},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.config)
			if tt.config != nil {
				tt.setup(builder)
				server := builder.Build()
				if tt.wantError && server != nil {
					t.Error("expected nil server")
				}
				if !tt.wantError && server == nil {
					t.Error("expected non-nil server")
				}
			} else if builder != nil {
				t.Error("expected nil builder with nil config")
			}
		})
	}
}
