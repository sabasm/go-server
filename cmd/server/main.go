package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/server"
	"github.com/sabasm/go-server/pkg/api/handlers/health"
	"github.com/sabasm/go-server/pkg/api/handlers/root"
	"github.com/sabasm/go-server/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	configLoader := config.NewConfigLoader()
	cfg, err := configLoader.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware(logger))

	srv := server.NewServerBuilder(cfg).
		WithRoute("/health", health.New().ServeHTTP).
		WithRoute("/", root.New().ServeHTTP).
		Build()

	serverError := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		handleServerShutdown(srv)
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	}
}

func handleServerShutdown(srv *server.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
