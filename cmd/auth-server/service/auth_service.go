package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sabasm/go-server/cmd/auth-server/handlers"
	"github.com/sabasm/go-server/internal/middleware"
	"go.uber.org/zap"
)

type AuthService struct {
	logger *zap.Logger
	router *mux.Router
}

func NewAuthService(logger *zap.Logger) *AuthService {
	return &AuthService{
		logger: logger,
		router: mux.NewRouter().StrictSlash(true),
	}
}

func (s *AuthService) ConfigureRoutes(m2mSecret string) {
	s.router.Use(middleware.LoggingMiddleware(s.logger))
	s.router.Use(middleware.RecoveryMiddleware(s.logger))

	api := s.router.PathPrefix("/auth").Subrouter()
	api.HandleFunc("/m2m/token", handlers.NewM2MHandler(m2mSecret).ServeHTTP).Methods(http.MethodPost)
	api.HandleFunc("/auth0/validate", handlers.NewAuth0Handler().ServeHTTP).Methods(http.MethodPost)
}

func (s *AuthService) GetHandler() http.Handler {
	return s.router
}
