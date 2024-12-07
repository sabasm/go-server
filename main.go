package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	env := os.Getenv("APP_ENV")
	port := os.Getenv("APP_PORT")

	fmt.Printf("Ejecutando en %s en el puerto %s\n", env, port)
}
