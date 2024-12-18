package server

import (
	"testing"
)

func TestGetAddress(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		expect string
	}{
		{
			name:   "Valid Config",
			config: &Config{Host: "localhost", Port: 8080},
			expect: "localhost:8080",
		},
		{
			name:   "Nil Config",
			config: nil,
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetAddress()
			if result != tt.expect {
				t.Errorf("expected %s, got %s", tt.expect, result)
			}
		})
	}
}
