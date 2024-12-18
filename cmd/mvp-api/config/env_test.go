package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "valid_env",
			envVars: map[string]string{
				"MVP_SECRET_KEY": "test-key",
				"APP_PORT":       "4004",
				"APP_HOST":       "localhost",
			},
			wantErr: false,
		},
		{
			name: "missing_secret_key",
			envVars: map[string]string{
				"APP_PORT": "4004",
				"APP_HOST": "localhost",
			},
			wantErr: true,
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

			_, err := LoadFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
