// cmd/mvp-api/main.go
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
	"github.com/sabasm/go-server/internal/api/core/handlers"
	"github.com/sabasm/go-server/internal/api/handlers/health"
	"github.com/sabasm/go-server/internal/api/handlers/root"
	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/middleware"
	"github.com/sabasm/go-server/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Configuration initialization and validation
	appConfig := config.LoadFromEnv()
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Logger setup with error handling
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger initialization failed: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Logger sync error: %v", err)
		}
	}()

	// Server configuration with validation
	srvCfg := &server.Config{
		Host:     appConfig.GetAppHost(),
		Port:     appConfig.GetAppPort(),
		BasePath: "/",
		Options: server.Options{
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	if err := srvCfg.Validate(); err != nil {
		logger.Fatal("Server configuration validation failed",
			zap.Error(err),
			zap.Any("config", srvCfg))
	}

	// Router setup with complete middleware chain
	router := mux.NewRouter()
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggingMiddleware(logger))

	// Server builder with all components
	srv := server.NewBuilder(srvCfg).
		WithLogger(logger).
		WithTimeout(
			srvCfg.Options.ReadTimeout,
			srvCfg.Options.WriteTimeout,
			srvCfg.Options.IdleTimeout,
		).
		WithRoute("/health", wrapHandler(health.New(), logger)).
		WithRoute("/", wrapHandler(root.New(), logger)).
		Build()

	// Server lifecycle management
	serverError := make(chan error, 1)
	go func() {
		logger.Info("Starting server",
			zap.String("host", srvCfg.Host),
			zap.Int("port", srvCfg.Port))
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error occurred", zap.Error(err))
			serverError <- err
		}
	}()

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		handleServerShutdown(srv, logger)
	case err := <-serverError:
		logger.Error("Fatal server error", zap.Error(err))
	}
}

func wrapHandler(h http.Handler, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bw := handlers.NewBufferedResponseWriter(w)
		defer func() {
			if err := bw.Flush(); err != nil {
				logger.Error("Response flush failed", zap.Error(err))
			}
		}()
		h.ServeHTTP(bw, r)
	}
}

func handleServerShutdown(srv server.ServerInterface, logger *zap.Logger) {
	logger.Info("Initiating graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Forced shutdown required", zap.Error(err))
	}
}
