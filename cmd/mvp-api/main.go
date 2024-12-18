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
	"github.com/sabasm/go-server/internal/api/routes"
	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/logger"
	"github.com/sabasm/go-server/internal/middleware"
	"github.com/sabasm/go-server/internal/server"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.LoadFromEnv()

	logr, err := logger.NewLogger([]string{"stdout"})
	if err != nil {
		log.Fatalf("Logger initialization failed: %v", err)
	}
	defer func() {
		if err := logr.Sync(); err != nil {
			log.Printf("Logger sync error: %v", err)
		}
	}()

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
		logr.Fatal("Server configuration validation failed", zap.Error(err), zap.Any("config", srvCfg))
	}

	router := mux.NewRouter()
	router.Use(middleware.RecoveryMiddleware(logr))
	router.Use(middleware.LoggingMiddleware(logr))

	defaultRoutes := routes.GetDefaultRoutes()
	routes.RegisterRoutes(router, defaultRoutes, func(h http.Handler) http.Handler {
		return wrapHandler(h, logr)
	})

	srv := server.NewBuilder(srvCfg).
		WithLogger(logr).
		WithTimeout(
			srvCfg.Options.ReadTimeout,
			srvCfg.Options.WriteTimeout,
			srvCfg.Options.IdleTimeout,
		).
		Build()

	srv.(*server.Server).SetHandler(router)

	serverError := make(chan error, 1)
	go func() {
		logr.Info("Starting server", zap.String("host", srvCfg.Host), zap.Int("port", srvCfg.Port))
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logr.Error("Server error occurred", zap.Error(err))
			serverError <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		handleServerShutdown(srv, logr)
	case err := <-serverError:
		logr.Error("Fatal server error", zap.Error(err))
	}
}

func wrapHandler(h http.Handler, logr *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bw := handlers.NewBufferedResponseWriter(w)
		defer func() {
			if err := bw.Flush(); err != nil {
				logr.Error("Response flush failed", zap.Error(err))
			}
		}()
		h.ServeHTTP(bw, r)
	})
}

func handleServerShutdown(srv server.ServerInterface, logr *zap.Logger) {
	logr.Info("Initiating graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logr.Error("Forced shutdown required", zap.Error(err))
	}
}
