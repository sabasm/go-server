package config

import (
	"fmt"

	baseconfig "github.com/sabasm/go-server/internal/config"
)

type MVPConfig struct {
	baseconfig.BaseConfig
	SecretKey string
}

func (c *MVPConfig) Validate() error {
	if err := c.BaseConfig.Validate(); err != nil {
		return err
	}
	if c.SecretKey == "" {
		return fmt.Errorf("secret key is required")
	}
	return nil
}

type Builder struct {
	config *MVPConfig
}

func NewBuilder() *Builder {
	return &Builder{
		config: &MVPConfig{
			BaseConfig: baseconfig.BaseConfig{
				AppPort: 4004,
				AppHost: "0.0.0.0",
			},
		},
	}
}

func (b *Builder) WithBaseConfig(base baseconfig.BaseConfig) *Builder {
	b.config.BaseConfig = base
	return b
}

func (b *Builder) WithSecretKey(key string) *Builder {
	b.config.SecretKey = key
	return b
}

func (b *Builder) Build() (*MVPConfig, error) {
	if err := b.config.Validate(); err != nil {
		return nil, err
	}
	return b.config, nil
}
