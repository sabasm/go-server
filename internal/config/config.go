package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Environment   string
	Port          int
	BaseURL       string
	RetryCount    int
	RetryDelay    int
	Monitoring    bool
	MetricsPrefix string
}

type ConfigLoader interface {
	LoadConfig() (*AppConfig, error)
}

type configLoader struct{}

func NewConfigLoader() ConfigLoader {
	return &configLoader{}
}

func (c *configLoader) LoadConfig() (*AppConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	config := &AppConfig{
		Environment:   getEnv("APP_ENV", "development"),
		Port:          getEnvAsInt("APP_PORT", 8080),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		RetryCount:    getEnvAsInt("RETRY_COUNT", 3),
		RetryDelay:    getEnvAsInt("RETRY_DELAY", 1000),
		Monitoring:    getEnvAsBool("MONITORING_ENABLED", false),
		MetricsPrefix: getEnv("METRICS_PREFIX", "app"),
	}

	return config, nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
