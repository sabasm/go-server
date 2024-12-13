// cmd/server/main.go
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
	"github.com/sabasm/go-server/internal/api/handlers/health"
	"github.com/sabasm/go-server/internal/api/handlers/root"
	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/middleware"
	"github.com/sabasm/go-server/internal/server"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.LoadFromEnv()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			log.Printf("Logger sync error: %v", syncErr)
		}
	}()

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware(logger))

	srvCfg := server.Config{
		Host:     appConfig.GetAppHost(),
		Port:     appConfig.GetAppPort(),
		BasePath: "/",
		Options: server.Options{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	srv := server.NewBuilder(&srvCfg).
		WithLogger(logger).
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

func handleServerShutdown(srv server.ServerInterface) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
