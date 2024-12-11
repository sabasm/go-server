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
