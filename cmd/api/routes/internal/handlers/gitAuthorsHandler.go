// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"net/http"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// func FetchDistinctAuthors() ([]string, error) {
// 	cursor, err := database.GitMetricsCollection.Distinct(context.Background(), "commited_by", bson.D{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var authors []string
// 	for _, author := range cursor {
// 		authors = append(authors, author.(string))
// 	}
// 	return authors, nil
// }

// // Git Repositories Handler
// func GitAuthorsHandler(w http.ResponseWriter, r *http.Request) {
// 	authors, err := FetchDistinctAuthors()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	data := models.GitMetricsViewData{
// 		Metrics: nil,
// 		Repos:   nil,
// 		Authors: authors,
// 	}

// 	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

//testing working but f() not
// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"net/http"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type GitMetricsCollectionInterface interface {
// 	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
// }

// var GitMetricsCollection GitMetricsCollectionInterface = database.GitMetricsCollection

// func FetchDistinctAuthors() ([]string, error) {
// 	cursor, err := GitMetricsCollection.Distinct(context.Background(), "commited_by", bson.D{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var authors []string
// 	for _, author := range cursor {
// 		authors = append(authors, author.(string))
// 	}
// 	return authors, nil
// }

// // Git Repositories Handler
// func GitAuthorsHandler(w http.ResponseWriter, r *http.Request) {
// 	authors, err := FetchDistinctAuthors()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	data := models.GitMetricsViewData{
// 		Metrics: nil,
// 		Repos:   nil,
// 		Authors: authors,
// 	}

// 	if err := helpers.RenderTemplateFunc(w, data, "git_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GitMetricsCollectionInterface interface {
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
}

var GitMetricsCollection GitMetricsCollectionInterface

func FetchDistinctAuthors() ([]string, error) {
	if GitMetricsCollection == nil {
		return nil, errors.New("GitMetricsCollection is not initialized")
	}

	cursor, err := GitMetricsCollection.Distinct(context.Background(), "commited_by", bson.D{})
	if err != nil {
		return nil, err
	}

	var authors []string
	for _, author := range cursor {
		if authorStr, ok := author.(string); ok {
			authors = append(authors, authorStr)
		} else {
			return nil, errors.New("type assertion failed for author")
		}
	}
	return authors, nil
}

// Git Repositories Handler
func GitAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := FetchDistinctAuthors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.GitMetricsViewData{
		Metrics: nil,
		Repos:   nil,
		Authors: authors,
	}

	if err := helpers.RenderTemplateFunc(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
