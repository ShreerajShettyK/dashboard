package handlers

import (
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"
)

// Git Repositories Handler
func GitHomeHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := FetchDistinctAuthors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repos, err := FetchDistinctRepositories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := models.GitMetricsViewData{
		Metrics: nil,
		Repos:   repos,
		Authors: authors,
	}

	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
