package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	router *mux.Router
	logger *zap.Logger
	srv    *http.Server
}

func (s *Server) Start() error {
	s.logger.Info("starting server",
		zap.String("addr", s.srv.Addr),
		zap.String("base_path", s.config.BasePath))
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.srv.Shutdown(ctx)
}
