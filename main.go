package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hossokawa/go-nethttp-example/internal/routes"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func run() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	connStr := os.Getenv("DB_URL")

	log.Println("Connecting to the database...")

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected to the database")

	router := http.NewServeMux()
	routes.SetupRoutes(router, conn)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on port :8080")
	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error starting the application: %s", err)
	}
}
