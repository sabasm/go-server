package logger

import (
	"os"
	"runtime"
	"testing"

	"go.uber.org/zap"
)

func TestConfigureLogger(t *testing.T) {
	tests := []struct {
		name       string
		appEnv     string
		expectLogs []string
	}{
		{
			name:   "Development Environment",
			appEnv: "development",
			expectLogs: []string{
				"stdout",
				"logs/development.log",
			},
		},
		{
			name:   "Production Environment",
			appEnv: "production",
			expectLogs: []string{
				"stdout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tt.appEnv)
			defer os.Unsetenv("APP_ENV")

			logger, err := ConfigureLogger()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if logger == nil {
				t.Fatal("expected non-nil logger")
			}
		})
	}
}

func TestSafeSync(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	if runtime.GOOS == "windows" {
		SafeSync(logger)
		t.Log("SafeSync should not fail on Windows")
	} else {
		SafeSync(logger)
		t.Log("SafeSync should complete without errors")
	}
}
