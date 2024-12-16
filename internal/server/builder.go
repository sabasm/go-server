package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type serverBuilder struct {
	config   *Config
	router   *mux.Router
	logger   *zap.Logger
	timeouts struct {
		read  time.Duration
		write time.Duration
		idle  time.Duration
	}
	middleware []mux.MiddlewareFunc
}

func NewBuilder(cfg *Config) *serverBuilder {
	if cfg == nil {
		return nil
	}
	return &serverBuilder{
		config:     cfg,
		router:     mux.NewRouter(),
		middleware: make([]mux.MiddlewareFunc, 0),
	}
}

func (b *serverBuilder) WithRoute(pattern string, handler http.HandlerFunc) *serverBuilder {
	if b == nil || b.router == nil {
		return b
	}
	b.router.HandleFunc(pattern, handler)
	return b
}

func (b *serverBuilder) WithMiddleware(m mux.MiddlewareFunc) *serverBuilder {
	if b == nil {
		return b
	}
	b.middleware = append(b.middleware, m)
	return b
}

func (b *serverBuilder) WithLogger(logger *zap.Logger) *serverBuilder {
	if b == nil {
		return b
	}
	b.logger = logger
	return b
}

func (b *serverBuilder) WithTimeout(read, write, idle time.Duration) *serverBuilder {
	if b == nil {
		return b
	}
	b.timeouts.read = read
	b.timeouts.write = write
	b.timeouts.idle = idle
	return b
}

func (b *serverBuilder) Build() ServerInterface {
	if b == nil || b.router == nil {
		return nil
	}

	if b.logger == nil {
		logger, _ := zap.NewProduction()
		b.logger = logger
	}

	for _, m := range b.middleware {
		b.router.Use(m)
	}

	return &Server{
		config: b.config,
		router: b.router,
		logger: b.logger,
		srv: &http.Server{
			Addr:         b.config.GetAddress(),
			Handler:      b.router,
			ReadTimeout:  b.timeouts.read,
			WriteTimeout: b.timeouts.write,
			IdleTimeout:  b.timeouts.idle,
		},
	}
}
