package main

import (
	"log"
	"net/http"

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
		WithRoute("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Service running")); err != nil {
				log.Printf("Failed to write response: %v", err)
			}
		}).
		Build()

	log.Printf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}


