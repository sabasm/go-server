package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sabasm/go-server/internal/api/handlers/health"
	"github.com/sabasm/go-server/internal/api/handlers/root"
)

type Route struct {
	Path    string
	Handler http.Handler
}

func GetDefaultRoutes() []Route {
	return []Route{
		{
			Path:    "/health",
			Handler: health.New(),
		},
		{
			Path:    "/",
			Handler: root.New(),
		},
	}
}

func RegisterRoutes(router *mux.Router, routes []Route, wrapHandler func(http.Handler) http.Handler) {
	for _, route := range routes {
		router.Handle(route.Path, wrapHandler(route.Handler))
	}
}
