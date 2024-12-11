package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *AppConfig
	}{
		{
			name:    "default_values",
			envVars: map[string]string{},
			expected: &AppConfig{
				Environment:   "development",
				Port:          8080,
				BaseURL:       "http://localhost:8080",
				RetryCount:    3,
				RetryDelay:    1000,
				Monitoring:    false,
				MetricsPrefix: "app",
			},
		},
		{
			name: "custom_values",
			envVars: map[string]string{
				"APP_ENV":            "production",
				"APP_PORT":           "9090",
				"BASE_URL":           "https://api.example.com",
				"RETRY_COUNT":        "5",
				"RETRY_DELAY":        "2000",
				"MONITORING_ENABLED": "true",
				"METRICS_PREFIX":     "prod",
			},
			expected: &AppConfig{
				Environment:   "production",
				Port:          9090,
				BaseURL:       "https://api.example.com",
				RetryCount:    5,
				RetryDelay:    2000,
				Monitoring:    true,
				MetricsPrefix: "prod",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// Test
			loader := NewConfigLoader()
			config, err := loader.LoadConfig()

			// Verify
			if err != nil {
				t.Errorf("LoadConfig() error = %v", err)
				return
			}

			if config.Environment != tt.expected.Environment {
				t.Errorf("Environment = %v, want %v", config.Environment, tt.expected.Environment)
			}
			if config.Port != tt.expected.Port {
				t.Errorf("Port = %v, want %v", config.Port, tt.expected.Port)
			}
			if config.BaseURL != tt.expected.BaseURL {
				t.Errorf("BaseURL = %v, want %v", config.BaseURL, tt.expected.BaseURL)
			}
			if config.RetryCount != tt.expected.RetryCount {
				t.Errorf("RetryCount = %v, want %v", config.RetryCount, tt.expected.RetryCount)
			}
			if config.RetryDelay != tt.expected.RetryDelay {
				t.Errorf("RetryDelay = %v, want %v", config.RetryDelay, tt.expected.RetryDelay)
			}
			if config.Monitoring != tt.expected.Monitoring {
				t.Errorf("Monitoring = %v, want %v", config.Monitoring, tt.expected.Monitoring)
			}
			if config.MetricsPrefix != tt.expected.MetricsPrefix {
				t.Errorf("MetricsPrefix = %v, want %v", config.MetricsPrefix, tt.expected.MetricsPrefix)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		setEnv     bool
		defaultVal string
		expected   string
	}{
		{
			name:       "existing_env_var",
			key:        "TEST_KEY",
			envValue:   "test_value",
			setEnv:     true,
			defaultVal: "default",
			expected:   "test_value",
		},
		{
			name:       "missing_env_var",
			key:        "TEST_KEY",
			setEnv:     false,
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "empty_env_var",
			key:        "TEST_KEY",
			envValue:   "",
			setEnv:     true,
			defaultVal: "default",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("getEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}
