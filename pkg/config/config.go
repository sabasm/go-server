package config

import "os"

type Config struct {
	AppHost string
	AppPort string
	Debug   bool
}

func LoadConfig() *Config {
	return &Config{
		AppHost: getEnv("APP_HOST", "localhost"),
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

func getEnvAsBool(name string, fallback bool) bool {
	if value, exists := os.LookupEnv(name); exists {
		return value == "true"
	}
	return fallback
}
