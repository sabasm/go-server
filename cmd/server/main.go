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
	"github.com/sabasm/go-server/internal/logger"
	"github.com/sabasm/go-server/internal/middleware"
	"github.com/sabasm/go-server/internal/server"
)

func main() {
	appConfig := config.LoadFromEnv()

	// Inicializar Logger
	logr, err := logger.NewLogger([]string{"stdout"})
	if err != nil {
		log.Fatalf("Error al inicializar el logger: %v", err)
	}
	defer func() {
		if syncErr := logr.Sync(); syncErr != nil {
			log.Printf("Error de sincronización del logger: %v", syncErr)
		}
	}()

	// Configurar Rutas
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware(logr))

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
		WithLogger(logr).
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
		log.Printf("Error del servidor: %v", err)
	}
}

func handleServerShutdown(srv server.ServerInterface) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Apagado forzado del servidor: %v", err)
	}
}
