package logger

import (
	"os"
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
			if tt.setTestMode {
				os.Setenv("LOGGER_MODE", "test")
				defer os.Unsetenv("LOGGER_MODE")
			}

			logger, err := NewLogger()

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
				if err != nil && err.Error() != "sync /dev/stderr: invalid argument" {
					t.Errorf("unexpected sync error: %v", err)
				}
			}()
		})
	}
}
