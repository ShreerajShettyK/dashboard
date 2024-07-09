package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchGitMetrics(userName, repoName string) ([]models.GitMetric, error) {
	filter := bson.M{}
	if userName != "None" {
		filter["commited_by"] = userName
	}
	if repoName != "None" {
		filter["reponame"] = repoName
	}
	opts := options.Find().SetSort(bson.D{{Key: "commit_date", Value: -1}})

	cursor, err := database.GitMetricsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	return helpers.DecodeCursor[models.GitMetric](context.Background(), cursor)
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
		Metrics:           metrics,
		Repos:             nil,
		Authors:           nil,
		RepoNameParameter: repoName,
		UserNameParameter: userName,
	}

	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
