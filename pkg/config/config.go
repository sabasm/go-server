package config

import (
	"os"
	"strconv"
)

// Config holds the application's configuration
type Config struct {
	AppPort string
	Debug   bool
}

// LoadConfig reads environment variables into Config
func LoadConfig() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		Debug:   getEnvAsBool("DEBUG", false),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return fallback
}
