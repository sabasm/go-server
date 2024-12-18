package config

import (
	"os"
	"testing"
)

const localhost = "localhost"

func TestLoadFromEnv(t *testing.T) {
	t.Run("valid env", func(t *testing.T) {
		os.Setenv("APP_HOST", "env_host")
		os.Setenv("APP_PORT", "9090")
		os.Setenv("DEBUG", "true")
		defer os.Clearenv()

		cfg := LoadFromEnv()
		if cfg.GetAppHost() != "env_host" {
			t.Errorf("expected AppHost 'env_host', got '%s'", cfg.GetAppHost())
		}
		if cfg.GetAppPort() != 9090 {
			t.Errorf("expected AppPort 9090, got %d", cfg.GetAppPort())
		}
		if !cfg.IsDebug() {
			t.Errorf("expected Debug to be true, got false")
		}
	})

	t.Run("missing env vars", func(t *testing.T) {
		os.Clearenv()
		cfg := LoadFromEnv()
		if cfg.GetAppHost() != localhost {
			t.Errorf("expected default AppHost 'localhost', got '%s'", cfg.GetAppHost())
		}
		if cfg.GetAppPort() != 8080 {
			t.Errorf("expected default AppPort 8080, got %d", cfg.GetAppPort())
		}
		if cfg.IsDebug() {
			t.Errorf("expected default Debug to be false, got true")
		}
	})
}

func TestEnvHelpers(t *testing.T) {
	t.Run("getEnv", func(t *testing.T) {
		os.Setenv("EXISTING_KEY", "value")
		defer os.Clearenv()
		if v := getEnv("EXISTING_KEY", "default"); v != "value" {
			t.Errorf("expected 'value', got '%s'", v)
		}
		if v := getEnv("MISSING_KEY", "default"); v != "default" {
			t.Errorf("expected 'default', got '%s'", v)
		}
	})

	t.Run("getEnvAsInt", func(t *testing.T) {
		os.Setenv("INT_KEY", "42")
		defer os.Clearenv()
		if v := getEnvAsInt("INT_KEY", 0); v != 42 {
			t.Errorf("expected 42, got %d", v)
		}
		if v := getEnvAsInt("MISSING_INT_KEY", 10); v != 10 {
			t.Errorf("expected 10, got %d", v)
		}
	})

	t.Run("getEnvAsBool", func(t *testing.T) {
		os.Setenv("BOOL_KEY", "true")
		defer os.Clearenv()
		if v := getEnvAsBool("BOOL_KEY", false); !v {
			t.Errorf("expected true, got false")
		}
		if v := getEnvAsBool("MISSING_BOOL_KEY", true); !v {
			t.Errorf("expected true, got false")
		}
	})
}
