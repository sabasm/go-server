package server

import (
	"context"
	"fmt"
	"hello-world-go/internal/config"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.AppConfig
	Router *mux.Router
	Logger *zap.Logger
	srv    *http.Server
}

type ServerBuilder struct {
	Config *config.AppConfig
	Router *mux.Router
	Logger *zap.Logger
}

func NewServerBuilder(cfg *config.AppConfig) *ServerBuilder {
	logger, _ := zap.NewProduction()
	return &ServerBuilder{
		Config: cfg,
		Router: mux.NewRouter(),
		Logger: logger,
	}
}

func (b *ServerBuilder) WithRoute(pattern string, handler http.HandlerFunc) *ServerBuilder {
	b.Router.HandleFunc(pattern, handler)
	return b
}

func (b *ServerBuilder) Build() *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", b.Config.Port),
		Handler:      b.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		Config: b.Config,
		Router: b.Router,
		Logger: b.Logger,
		srv:    srv,
	}
}

func (s *Server) Start() error {
	s.Logger.Info("Starting server", zap.String("addr", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Server shutting down")
	return s.srv.Shutdown(ctx)
}
