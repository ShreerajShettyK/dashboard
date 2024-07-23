package routes

import (
	"dashboard/cmd/api/routes/internal/handlerAwsEfficiency"
	"dashboard/cmd/api/routes/internal/handlers"
	"dashboard/cmd/api/routes/internal/handlersAws"
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	Methods(methods ...string) *mux.Route
}

func RegisterRoutes(r Router) {
	r.HandleFunc("/aws_metrics/home", handlersAws.AwsHomeHandler).Methods("GET")
	r.HandleFunc("/aws_metrics/home/resources", handlersAws.AWSMetricsHandler).Methods("GET")

	r.HandleFunc("/git_metrics/home", handlers.GitHomeHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/commits", handlers.GitMetricsHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/repos", handlers.GitReposHandler).Methods("GET")
	r.HandleFunc("/git_metrics/home/authors", handlers.GitAuthorsHandler).Methods("GET")
	r.HandleFunc("/git_metrics/repoAuthors", handlers.GitAuthorsByRepoHandler).Methods("GET")

	// r.HandleFunc("/aws_billing/instances", handlerAwsEfficiency.ListInstancesHandler).Methods("GET")
	// r.HandleFunc("/aws_billing/instance/{instance_id}", handlerAwsEfficiency.InstanceDetailsHandler).Methods("GET")

	r.HandleFunc("/aws_billing/services", handlerAwsEfficiency.ListServicesHandler).Methods("GET")
	r.HandleFunc("/aws_billing/service/{service}/instances", handlerAwsEfficiency.ListServiceInstancesHandler).Methods("GET")
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	RegisterRoutes(r)
	return r
}
