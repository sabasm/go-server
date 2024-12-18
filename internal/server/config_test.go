package server

import (
	"testing"
	"time"
)

const localhost = "localhost"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{"Valid Config", &Config{Host: localhost, Port: 8080, Options: Options{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}}, false},
		{"Invalid Port", &Config{Host: localhost, Port: 70000}, true},
		{"Empty Host", &Config{Host: "", Port: 8080}, true},
		{"Negative Read Timeout", &Config{Host: localhost, Port: 8080, Options: Options{
			ReadTimeout: -5 * time.Second,
		}}, true},
		{"Negative Write Timeout", &Config{Host: localhost, Port: 8080, Options: Options{
			WriteTimeout: -5 * time.Second,
		}}, true},
		{"Negative Idle Timeout", &Config{Host: localhost, Port: 8080, Options: Options{
			IdleTimeout: -5 * time.Second,
		}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
