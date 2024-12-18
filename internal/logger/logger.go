package logger

import (
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const windowsOS = "windows"

type LoggerOptions struct {
	OutputPaths []string
}

func NewLogger(outputs []string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	if os.Getenv("DEBUG") == "true" {
		level.SetLevel(zap.DebugLevel)
	}

	config := zap.Config{
		Encoding:         "json",
		Level:            level,
		OutputPaths:      outputs,
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			LevelKey:     "level",
			MessageKey:   "msg",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	return config.Build()
}

func ConfigureLogger() (*zap.Logger, error) {
	var outputPaths []string

	if os.Getenv("APP_ENV") == "development" {
		logDir := "logs"
		logFile := filepath.Join(logDir, "development.log")

		// Ensure the logs directory exists
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}
		outputPaths = []string{"stdout", logFile}
	} else {
		outputPaths = []string{"stdout"}
	}

	return NewLogger(outputPaths)
}

// SafeSync safely syncs the logger
func SafeSync(logger *zap.Logger) {
	if runtime.GOOS != windowsOS {
		_ = logger.Sync()
	}
}
