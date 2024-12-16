package server

import (
	"fmt"
	"time"
)

type Config struct {
	Host     string
	Port     int
	BasePath string
	Options  Options
}

type Options struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

func (c *Config) GetAddress() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Port)
	}
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Options.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}
	if c.Options.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}
	if c.Options.IdleTimeout <= 0 {
		return fmt.Errorf("idle timeout must be positive")
	}
	return nil
}
