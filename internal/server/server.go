package server

import (
	"fmt"
	"net/http"

	"hello-world-go/internal/config"
)

type ServerBuilder struct {
	Config *config.AppConfig
	Router *http.ServeMux
}

func NewServerBuilder(cfg *config.AppConfig) *ServerBuilder {
	return &ServerBuilder{
		Config: cfg,
		Router: http.NewServeMux(),
	}
}

func (b *ServerBuilder) WithRoute(pattern string, handler http.HandlerFunc) *ServerBuilder {
	b.Router.HandleFunc(pattern, handler)
	return b
}

func (b *ServerBuilder) Build() *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", b.Config.Port),
		Handler: b.Router,
	}
}


