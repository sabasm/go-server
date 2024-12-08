package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("BASE_URL", "http://test.local")
	os.Setenv("RETRY_COUNT", "5")
	os.Setenv("RETRY_DELAY", "2000")
	os.Setenv("MONITORING_ENABLED", "true")
	os.Setenv("METRICS_PREFIX", "testapp")

	configLoader := NewConfigLoader()
	cfg, err := configLoader.LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Environment != "test" || cfg.Port != 9090 {
		t.Errorf("Expected APP_ENV=test and APP_PORT=9090, got %v and %v", cfg.Environment, cfg.Port)
	}
}

func TestGetEnvAsIntInvalid(t *testing.T) {
	os.Setenv("APP_PORT", "x")
	val := getEnvAsInt("APP_PORT", 8080)
	if val != 8080 {
		t.Errorf("Expected 8080, got %d", val)
	}
}

func TestGetEnvAsBoolInvalid(t *testing.T) {
	os.Setenv("MONITORING_ENABLED", "x")
	val := getEnvAsBool("MONITORING_ENABLED", false)
	if val != false {
		t.Errorf("Expected false, got %v", val)
	}
}
