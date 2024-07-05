package main

import (
	"dashboard/cmd/api/routes/internal/handlers"
	"dashboard/cmd/api/routes/internal/middleware"
	"log"
	"net/http"
)

func main() {

	http.Handle("/send-message", middleware.AuthMiddleware(http.HandlerFunc(handlers.SendMessageHandler)))

	log.Println("Server starting at :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
