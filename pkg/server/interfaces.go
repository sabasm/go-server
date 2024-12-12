package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerInterface interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type ServerBuilder interface {
	WithRoute(pattern string, handler http.HandlerFunc) ServerBuilder
	WithMiddleware(middleware mux.MiddlewareFunc) ServerBuilder
	WithLogger(logger *zap.Logger) ServerBuilder
	WithTimeout(read, write, idle time.Duration) ServerBuilder
	Build() ServerInterface
}
