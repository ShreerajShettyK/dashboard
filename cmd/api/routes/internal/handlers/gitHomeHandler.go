package handlers

import (
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"
)

// Make these variables to allow mocking in tests
var FetchDistinctAuthorsMock = FetchDistinctAuthors
var FetchDistinctRepositoriesMock = FetchDistinctRepositories
var RenderTemplate = helpers.RenderTemplate

// Git Repositories Handler
func GitHomeHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := FetchDistinctAuthorsMock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repos, err := FetchDistinctRepositoriesMock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.GitMetricsViewData{
		Metrics: nil,
		Repos:   repos,
		Authors: authors,
	}

	if err := RenderTemplate(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
