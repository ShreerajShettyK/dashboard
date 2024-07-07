// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"html/template"
// 	"net/http"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func FetchGitMetrics(userName, repoName string) ([]models.GitMetric, error) {
// 	filter := bson.M{
// 		"user_name": userName,
// 		"repo_name": repoName,
// 	}
// 	opts := options.Find().SetSort(bson.D{{Key: "commit_date", Value: -1}})

// 	cursor, err := database.GitMetricsCollection.Find(context.Background(), filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	var metrics []models.GitMetric
// 	for cursor.Next(context.Background()) {
// 		var metric models.GitMetric
// 		if err := cursor.Decode(&metric); err != nil {
// 			return nil, err
// 		}
// 		metrics = append(metrics, metric)
// 	}
// 	return metrics, nil
// }

// func GitMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	userName := r.URL.Query().Get("user_name")
// 	repoName := r.URL.Query().Get("repo_name")

// 	metrics, err := FetchGitMetrics(userName, repoName)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Use the relative path
// 	tmplPath := "internal/templates/git_dashboard.html"

// 	tmpl, err := template.ParseFiles(tmplPath)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Execute the template with metrics data
// 	if err := tmpl.Execute(w, metrics); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"

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

	// Print the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory:", err)
	} else {
		log.Println("Current working directory:", cwd)
	}

	// Use the relative path
	tmplPath := "internal/templates/git_dashboard.html"
	log.Println("Template Path:", tmplPath) // Debug print

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with metrics data
	if err := tmpl.Execute(w, metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
