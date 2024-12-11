package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerBuilder interface {
	WithRoute(pattern string, handler http.HandlerFunc) ServerBuilder
	WithMiddleware(middleware mux.MiddlewareFunc) ServerBuilder
	WithLogger(logger *zap.Logger) ServerBuilder
	WithTimeout(read, write, idle time.Duration) ServerBuilder
	Build() *Server
}

type serverBuilder struct {
	config   *Config
	router   *mux.Router
	logger   *zap.Logger
	timeouts struct {
		read  time.Duration
		write time.Duration
		idle  time.Duration
	}
}

func NewBuilder(cfg *Config) ServerBuilder {
	return &serverBuilder{
		config: cfg,
		router: mux.NewRouter(),
	}
}

func (b *serverBuilder) WithRoute(pattern string, handler http.HandlerFunc) ServerBuilder {
	b.router.HandleFunc(pattern, handler)
	return b
}

func (b *serverBuilder) WithMiddleware(middleware mux.MiddlewareFunc) ServerBuilder {
	b.router.Use(middleware)
	return b
}

func (b *serverBuilder) WithLogger(logger *zap.Logger) ServerBuilder {
	b.logger = logger
	return b
}

func (b *serverBuilder) WithTimeout(read, write, idle time.Duration) ServerBuilder {
	b.timeouts.read = read
	b.timeouts.write = write
	b.timeouts.idle = idle
	return b
}

func (b *serverBuilder) Build() *Server {
	if b.logger == nil {
		b.logger, _ = zap.NewProduction()
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
