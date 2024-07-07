package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchGitMetrics(userName, repoName string) ([]models.GitMetric, error) {
	filter := bson.M{
		"commited_by": userName,
		"reponame":    repoName,
	}
	opts := options.Find().SetSort(bson.D{{Key: "commit_date", Value: -1}})

	cursor, err := database.GitMetricsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var metrics []models.GitMetric
	for cursor.Next(context.Background()) {
		var metric models.GitMetric
		if err := cursor.Decode(&metric); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func GitMetricsHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("user_name")
	repoName := r.URL.Query().Get("repo_name")

	metrics, err := FetchGitMetrics(userName, repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.GitMetricsViewData{
		Metrics: metrics,
		Repos:   nil,
		Authors: nil,
	}

	// Use the relative path
	tmplPath := "internal/templates/git_dashboard.html"
	// log.Println("Template Path:", tmplPath) // Debug print

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with metrics data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Fetched records")
}
