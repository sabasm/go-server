// internal/config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *BaseConfig
	}{
		{
			name:    "default_values",
			envVars: map[string]string{},
			expected: &BaseConfig{
				AppHost: "localhost",
				AppPort: 8080,
				Debug:   false,
			},
		},
		{
			name: "custom_values",
			envVars: map[string]string{
				"APP_HOST": "0.0.0.0",
				"APP_PORT": "3000",
				"DEBUG":    "true",
			},
			expected: &BaseConfig{
				AppHost: "0.0.0.0",
				AppPort: 3000,
				Debug:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			got := LoadFromEnv().(*BaseConfig)
			if got.AppHost != tt.expected.AppHost {
				t.Errorf("AppHost = %v, want %v", got.AppHost, tt.expected.AppHost)
			}
			if got.AppPort != tt.expected.AppPort {
				t.Errorf("AppPort = %v, want %v", got.AppPort, tt.expected.AppPort)
			}
			if got.Debug != tt.expected.Debug {
				t.Errorf("Debug = %v, want %v", got.Debug, tt.expected.Debug)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		setEnv   bool
		want     string
	}{
		{
			name:     "exists",
			key:      "TEST_KEY",
			value:    "test_value",
			fallback: "default",
			setEnv:   true,
			want:     "test_value",
		},
		{
			name:     "not_exists",
			key:      "TEST_KEY",
			fallback: "default",
			setEnv:   false,
			want:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := getEnv(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsBool(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback bool
		setEnv   bool
		want     bool
	}{
		{
			name:     "true_value",
			key:      "TEST_BOOL",
			value:    "true",
			fallback: false,
			setEnv:   true,
			want:     true,
		},
		{
			name:     "false_value",
			key:      "TEST_BOOL",
			value:    "false",
			fallback: true,
			setEnv:   true,
			want:     false,
		},
		{
			name:     "invalid_value",
			key:      "TEST_BOOL",
			value:    "invalid",
			fallback: false,
			setEnv:   true,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := getEnvAsBool(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnvAsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppHost(t *testing.T) {
	config := &BaseConfig{AppHost: "test_host"}
	if config.GetAppHost() != "test_host" {
		t.Errorf("Expected AppHost to be 'test_host', got '%v'", config.GetAppHost())
	}
}

func TestGetAppPort(t *testing.T) {
	config := &BaseConfig{AppPort: 1234}
	if config.GetAppPort() != 1234 {
		t.Errorf("Expected AppPort to be 1234, got %v", config.GetAppPort())
	}
}

func TestIsDebug(t *testing.T) {
	config := &BaseConfig{Debug: true}
	if !config.IsDebug() {
		t.Errorf("Expected IsDebug to be true, got false")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *BaseConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &BaseConfig{
				AppHost: "localhost",
				AppPort: 8080,
			},
			wantErr: false,
		},
		{
			name: "invalid_port",
			config: &BaseConfig{
				AppHost: "localhost",
				AppPort: 0,
			},
			wantErr: true,
			errMsg:  "expected positive port number",
		},
		{
			name: "invalid_host",
			config: &BaseConfig{
				AppHost: "",
				AppPort: 8080,
			},
			wantErr: true,
			errMsg:  "expected non-empty host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
