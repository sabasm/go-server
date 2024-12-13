package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sabasm/go-server/internal/config"
	"github.com/sabasm/go-server/internal/server"
	"github.com/sabasm/go-server/pkg/api/handlers/health"
	"github.com/sabasm/go-server/pkg/api/handlers/root"
)

func main() {
	configLoader := config.NewConfigLoader()
	appConfig, err := configLoader.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv := server.NewServerBuilder(appConfig).
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			cancel()
			log.Printf("Server forced to shutdown: %v", err)
			os.Exit(1)
		}
		cancel()
	case err := <-serverError:
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
