package main

import (
	"log"
	"net/http"

	"dashboard/cmd/api/routes/config"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/routes"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize the database connection
	database.InitDB(cfg.MongoDBURI, cfg.DBName)

	// Set up the router
	r := routes.RegisterRoutes()

	// Start the server
	log.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
