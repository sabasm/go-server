package config

import (
	"testing"

	baseconfig "github.com/sabasm/go-server/internal/config"
)

func TestMVPConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *MVPConfig
		wantErr bool
	}{
		{
			name: "valid_config",
			config: &MVPConfig{
				BaseConfig: baseconfig.BaseConfig{
					AppHost: "localhost",
					AppPort: 4004,
				},
				SecretKey: "test-key",
			},
			wantErr: false,
		},
		{
			name: "missing_secret_key",
			config: &MVPConfig{
				BaseConfig: baseconfig.BaseConfig{
					AppHost: "localhost",
					AppPort: 4004,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MVPConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
