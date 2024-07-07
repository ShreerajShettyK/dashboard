package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"log"
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
	log.Println("Fetched repositories")
}
