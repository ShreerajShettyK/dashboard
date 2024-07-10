package main

import (
	"log"
	"net/http"

	"dashboard/cmd/api/routes/config"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/handlers"
	"dashboard/cmd/api/routes/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	database.InitDB(cfg.MongoDBURI, cfg.DBName)
	handlers.GitMetricsCollection = database.GitMetricsCollection

	r := routes.RegisterRoutes()

	log.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
