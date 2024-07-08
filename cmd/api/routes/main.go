// package main

// import (
// 	"dashboard/cmd/api/routes/internal/handlers"
// 	"dashboard/cmd/api/routes/internal/middleware"
// 	"log"
// 	"net/http"
// )

// func main() {

// 	http.Handle("/send-message", middleware.AuthMiddleware(http.HandlerFunc(handlers.SendMessageHandler)))

// 	log.Println("Server starting at :8000")
// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatalf("Server error: %v", err)
// 	}
// }

package main

import (
	"log"
	"net/http"
	"os"

	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	if mongoURI == "" {
		log.Fatal("MONGODB connectionstring or db name not set")
	}

	// Initialize the database connection
	database.InitDB(mongoURI, dbName)

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/aws_metrics", handlers.AWSMetricsHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home", handlers.GitHomeHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/commits", handlers.GitMetricsHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/repos", handlers.GitReposHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/authors", handlers.GitAuthorsHandler).Methods("GET")

	// Start the server
	log.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
