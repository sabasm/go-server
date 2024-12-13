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
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Welcome to My MVP App"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func main() {
	appConfig := config.LoadFromEnv()
	port := appConfig.GetAppPort()

	cfg := &server.Config{
		Host:     appConfig.GetAppHost(),
		Port:     port,
		BasePath: "/",
		Options: server.Options{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	log.Printf("Loaded config: %+v", cfg)

	srv := server.NewBuilder(cfg).
		WithRoute("/", rootHandler).
		WithRoute("/health", healthHandler).
		Build()

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
		shutdownCancel()
		os.Exit(1)
	}
	shutdownCancel()
	log.Println("Server stopped gracefully")
}
