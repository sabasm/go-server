package config

import (
	"testing"
)

const (
	expectedPort        = 8080
	expectedHost        = "localhost"
	expectedLogLevel    = "info"
	expectedEnvironment = "test"
)

func TestOptions(t *testing.T) {
	opts := Options{}
	opts.Server.Port = expectedPort
	opts.Server.Host = expectedHost
	opts.Logger.Level = expectedLogLevel
	opts.App.Environment = expectedEnvironment

	if opts.Server.Port != expectedPort {
		t.Errorf("expected port %d, got %d", expectedPort, opts.Server.Port)
	}

	if opts.Logger.Level != expectedLogLevel {
		t.Errorf("expected level %s, got %s", expectedLogLevel, opts.Logger.Level)
	}

	if opts.App.Environment != expectedEnvironment {
		t.Errorf("expected env %s, got %s", expectedEnvironment, opts.App.Environment)
	}
}
