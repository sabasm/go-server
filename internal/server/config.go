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
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("expected positive port number")
	}
	if c.Host == "" {
		return fmt.Errorf("expected non-empty host")
	}
	return nil
}
