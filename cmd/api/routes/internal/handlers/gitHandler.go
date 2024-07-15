package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchGitMetrics(userName, repoName string, limit, skip int64) ([]models.GitMetric, error) {
	filter := bson.M{}
	if userName != "None" {
		filter["commited_by"] = userName
	}
	if repoName != "None" {
		filter["reponame"] = repoName
	}
	opts := options.Find().SetSort(bson.D{{Key: "commit_date", Value: -1}}).SetLimit(limit).SetSkip(skip)

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
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit := int64(10)
	skip := (page - 1) * limit

	metrics, err := FetchGitMetrics(userName, repoName, limit, skip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format commit dates
	for i := range metrics {
		metrics[i].FormattedCommitDate = metrics[i].CommitDate.Format("Jan 2, 2006 at 3:04pm")
	}

	// Add pagination details
	data := models.GitMetricsViewData{
		Metrics:           metrics,
		Repos:             nil,
		Authors:           nil,
		RepoNameParameter: repoName,
		UserNameParameter: userName,
		CurrentPage:       page,
		PreviousPage:      page - 1,
		NextPage:          page + 1,
	}

	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
