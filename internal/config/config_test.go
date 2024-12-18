package config

import (
	"testing"
)

const localhost = "localhost"

func TestGetAppHost(t *testing.T) {
	config := &BaseConfig{AppHost: localhost}
	if config.GetAppHost() != localhost {
		t.Errorf("expected %s, got %s", localhost, config.GetAppHost())
	}
}

func TestGetAppPort(t *testing.T) {
	config := &BaseConfig{AppPort: 8080}
	if config.GetAppPort() != 8080 {
		t.Errorf("expected 8080, got %d", config.GetAppPort())
	}
}

func TestIsDebug(t *testing.T) {
	config := &BaseConfig{Debug: true}
	if !config.IsDebug() {
		t.Errorf("expected true, got false")
	}
}

func TestValidate(t *testing.T) {
	validConfig := &BaseConfig{AppHost: localhost, AppPort: 8080}
	if err := validConfig.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}

	invalidConfig := &BaseConfig{AppHost: "", AppPort: -1}
	if err := invalidConfig.Validate(); err == nil {
		t.Error("expected validation error, got none")
	}
}
