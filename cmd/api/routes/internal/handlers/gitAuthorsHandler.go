package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func FetchDistinctAuthors() ([]string, error) {
	cursor, err := database.GitMetricsCollection.Distinct(context.Background(), "commited_by", bson.D{})
	if err != nil {
		return nil, err
	}

	var authors []string
	for _, author := range cursor {
		authors = append(authors, author.(string))
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

	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
