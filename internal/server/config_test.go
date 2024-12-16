package server

import (
	"testing"
	"time"
)

func TestConfigValidations(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "valid_config",
			config: &Config{
				Host: "localhost",
				Port: 8080,
				Options: Options{
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 5 * time.Second,
					IdleTimeout:  30 * time.Second,
				},
			},
			wantError: false,
		},
		{
			name:      "nil_config",
			config:    nil,
			wantError: true,
		},
		{
			name: "invalid_port",
			config: &Config{
				Host: "localhost",
				Port: -1,
			},
			wantError: true,
		},
		{
			name: "invalid_host",
			config: &Config{
				Host: "",
				Port: 8080,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
