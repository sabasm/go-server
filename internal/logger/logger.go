package logger

import (
	"os"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	if os.Getenv("LOGGER_MODE") == "test" {
		return zap.NewDevelopment(zap.WithCaller(false))
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stdout"}

	return config.Build()
}
