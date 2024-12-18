package config

import (
	"fmt"
	"os"

	baseconfig "github.com/sabasm/go-server/internal/config"
)

func LoadFromEnv() (baseconfig.Config, error) {
	base := baseconfig.LoadFromEnv()
	secretKey := os.Getenv("MVP_SECRET_KEY")

	if secretKey == "" {
		return nil, fmt.Errorf("MVP_SECRET_KEY environment variable is required")
	}

	config, err := NewBuilder().
		WithBaseConfig(*base.(*baseconfig.BaseConfig)).
		WithSecretKey(secretKey).
		Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	return config, nil
}
