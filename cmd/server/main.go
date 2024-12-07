package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hello-world-go/internal/config"
	"hello-world-go/internal/server"
)

func main() {
	configLoader := config.NewConfigLoader()
	appConfig, err := configLoader.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv := server.NewServerBuilder(appConfig).
		WithRoute("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				log.Printf("Failed to write health response: %v", err)
			}
		}).
		WithRoute("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Service running")); err != nil {
				log.Printf("Failed to write response: %v", err)
			}
		}).
		Build()

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
