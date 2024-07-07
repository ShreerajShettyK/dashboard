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
	log.Println("Fetched authors")
}
