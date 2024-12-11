package server

import (
	"fmt"
	"time"
)

type Config struct {
	Port     int
	Host     string
	BasePath string
	Options  struct {
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
