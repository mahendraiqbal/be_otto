package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mahendraiqbal/be_otto/internal/api"
	"github.com/mahendraiqbal/be_otto/internal/api/handlers"
	"github.com/mahendraiqbal/be_otto/internal/repository/postgres"
)

func main() {
	// Get database connection string from environment variable
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5433/voucher_db?sslmode=disable"
	}

	// Initialize repository
	repo, err := postgres.NewPostgresRepository(dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	handler := handlers.NewHandler(repo)

	// Setup routes
	router := api.SetupRoutes(handler)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
