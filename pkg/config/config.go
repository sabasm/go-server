package config

import (
	"os"
	"strconv"
	"sync"
)

type Config interface {
	GetAppHost() string
	GetAppPort() int
	IsDebug() bool
}

type BaseConfig struct {
	AppHost string
	AppPort int
	Debug   bool
}

func (c *BaseConfig) GetAppHost() string {
	return c.AppHost
}

func (c *BaseConfig) GetAppPort() int {
	return c.AppPort
}

func (c *BaseConfig) IsDebug() bool {
	return c.Debug
}

type ConfigBuilder struct {
	config *BaseConfig
	mutex  sync.Mutex
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &BaseConfig{},
	}
}

func (b *ConfigBuilder) WithAppHost(host string) *ConfigBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.config.AppHost = host
	return b
}

func (b *ConfigBuilder) WithAppPort(port int) *ConfigBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.config.AppPort = port
	return b
}

func (b *ConfigBuilder) WithDebug(debug bool) *ConfigBuilder {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.config.Debug = debug
	return b
}

func (b *ConfigBuilder) Build() Config {
	return b.config
}

func LoadFromEnv() Config {
	builder := NewConfigBuilder()

	appHost := getEnv("APP_HOST", "localhost")
	appPort := getEnvAsInt("APP_PORT", 8080)
	debug := getEnvAsBool("DEBUG", false)

	return builder.
		WithAppHost(appHost).
		WithAppPort(appPort).
		WithDebug(debug).
		Build()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return fallback
}
