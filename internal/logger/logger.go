package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOptions struct {
	OutputPaths []string
}

func NewLogger(outputs []string) (*zap.Logger, error) {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
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
