package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func FetchDistinctRepositories() ([]string, error) {
	cursor, err := database.GitMetricsCollection.Distinct(context.Background(), "reponame", bson.D{})
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, repo := range cursor {
		repos = append(repos, repo.(string))
	}
	return repos, nil
}

// Git Repositories Handler
func GitReposHandler(w http.ResponseWriter, r *http.Request) {
	repos, err := FetchDistinctRepositories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.GitMetricsViewData{
		Metrics: nil,
		Repos:   repos,
		Authors: nil,
	}

	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
