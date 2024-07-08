package handlers

import (
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
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

	// Use the relative path
	tmplPath := "internal/templates/git_dashboard.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with repos data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
