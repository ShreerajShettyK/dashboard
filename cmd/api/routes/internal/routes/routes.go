package routes

import (
	"dashboard/cmd/api/routes/internal/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/aws_metrics/home", handlers.AwsHomeHandler).Methods("GET")
	r.HandleFunc("/aws_metrics/home/resources", handlers.AWSMetricsHandler).Methods("GET")

	r.HandleFunc("/git_metrics/home", handlers.GitHomeHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/commits", handlers.GitMetricsHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/repos", handlers.GitReposHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/authors", handlers.GitAuthorsHandler).Methods("GET")

	return r
}
