package logger

import (
	"os"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		setTestMode   bool
		expectedError bool
	}{
		{
			name:          "Production Logger",
			setTestMode:   false,
			expectedError: false,
		},
		{
			name:          "Test Logger",
			setTestMode:   true,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outputPaths []string
			if tt.setTestMode {
				// Redirect logger output to /dev/null during tests to avoid sync errors
				outputPaths = []string{os.DevNull}
				os.Setenv("LOG_LEVEL", "debug") // Example: set debug level if needed
				defer os.Unsetenv("LOG_LEVEL")
			}

			logger, err := NewLogger(outputPaths)

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if logger == nil {
				t.Error("expected non-nil logger")
			}

			defer func() {
				err := logger.Sync()
				if err != nil && !strings.Contains(err.Error(), "invalid argument") {
					t.Errorf("unexpected sync error: %v", err)
				}
				// Alternatively, you can choose to ignore specific errors
			}()
		})
	}
}
