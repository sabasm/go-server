package config

import "testing"

const (
	expectedLogLevel    = "info"
	expectedEnvironment = "test"
)

func TestOptions(t *testing.T) {
	opts := Options{}
	opts.Server.Port = 8080
	opts.Server.Host = "localhost"
	opts.Logger.Level = expectedLogLevel
	opts.App.Environment = expectedEnvironment

	if opts.Server.Port != 8080 {
		t.Errorf("expected port %d, got %d", 8080, opts.Server.Port)
	}

	if opts.Logger.Level != expectedLogLevel {
		t.Errorf("expected level %s, got %s", expectedLogLevel, opts.Logger.Level)
	}

	if opts.App.Environment != expectedEnvironment {
		t.Errorf("expected env %s, got %s", expectedEnvironment, opts.App.Environment)
	}
}
