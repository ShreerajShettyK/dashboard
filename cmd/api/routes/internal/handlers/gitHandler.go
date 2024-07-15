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

// testing:
// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"net/http"
// 	"strconv"

// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type NewGitMetricsCollectionInterface interface {
// 	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error)
// }

// // MockGitMetricsCollection is a mock implementation of GitMetricsCollectionInterface
// type NewMockGitMetricsCollection struct {
// 	mock.Mock
// }

// func (m *NewMockGitMetricsCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error) {
// 	args := m.Called(ctx, filter, opts)
// 	return args.Get(0).(CursorInterface), args.Error(1)
// }

// type CursorInterface interface {
// 	Close(ctx context.Context) error
// 	All(ctx context.Context, results interface{}) error
// }

// var NewGitMetricsCollection NewGitMetricsCollectionInterface

// func FetchGitMetrics(userName, repoName string, limit, skip int64) ([]models.GitMetric, error) {
// 	filter := bson.M{}
// 	if userName != "None" {
// 		filter["commited_by"] = userName
// 	}
// 	if repoName != "None" {
// 		filter["reponame"] = repoName
// 	}
// 	opts := options.Find().SetSort(bson.D{{Key: "commit_date", Value: -1}}).SetLimit(limit).SetSkip(skip)

// 	cursor, err := NewGitMetricsCollection.Find(context.Background(), filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	var metrics []models.GitMetric
// 	if err := cursor.All(context.Background(), &metrics); err != nil {
// 		return nil, err
// 	}

// 	return metrics, nil
// }

// func GitMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	userName := r.URL.Query().Get("user_name")
// 	repoName := r.URL.Query().Get("repo_name")
// 	pageStr := r.URL.Query().Get("page")
// 	page, err := strconv.ParseInt(pageStr, 10, 64)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	limit := int64(10)
// 	skip := (page - 1) * limit

// 	metrics, err := FetchGitMetrics(userName, repoName, limit, skip)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Format commit dates
// 	for i := range metrics {
// 		metrics[i].FormattedCommitDate = metrics[i].CommitDate.Format("Jan 2, 2006 at 3:04pm")
// 	}

// 	// Add pagination details
// 	data := models.GitMetricsViewData{
// 		Metrics:           metrics,
// 		Repos:             nil,
// 		Authors:           nil,
// 		RepoNameParameter: repoName,
// 		UserNameParameter: userName,
// 		CurrentPage:       page,
// 		PreviousPage:      page - 1,
// 		NextPage:          page + 1,
// 	}

// 	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

//below code testing 65% im getting but page doesnt load (functionality fails)

// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"net/http"
// 	"strconv"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // NewGitMetricsCollectionInterface defines the methods we need from the MongoDB collection
// type NewGitMetricsCollectionInterface interface {
// 	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
// }

// // CursorInterface defines the methods we need from the MongoDB cursor
// type CursorInterface interface {
// 	Close(ctx context.Context) error
// 	All(ctx context.Context, results interface{}) error
// }

// var NewGitMetricsCollection NewGitMetricsCollectionInterface

// func FetchGitMetrics(userName, repoName string, limit, skip int64) ([]models.GitMetric, error) {
// 	filter := bson.M{}
// 	if userName != "None" {
// 		filter["commited_by"] = userName
// 	}
// 	if repoName != "None" {
// 		filter["reponame"] = repoName
// 	}
// 	opts := options.Find().SetSort(bson.D{{Key: "commit_date", Value: -1}}).SetLimit(limit).SetSkip(skip)

// 	cursor, err := NewGitMetricsCollection.Find(context.Background(), filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	var metrics []models.GitMetric
// 	if err := cursor.All(context.Background(), &metrics); err != nil {
// 		return nil, err
// 	}
// 	return metrics, nil
// }

// func GitMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	userName := r.URL.Query().Get("user_name")
// 	repoName := r.URL.Query().Get("repo_name")
// 	pageStr := r.URL.Query().Get("page")
// 	page, err := strconv.ParseInt(pageStr, 10, 64)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	limit := int64(10)
// 	skip := (page - 1) * limit

// 	metrics, err := FetchGitMetrics(userName, repoName, limit, skip)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Format commit dates
// 	for i := range metrics {
// 		metrics[i].FormattedCommitDate = metrics[i].CommitDate.Format("Jan 2, 2006 at 3:04pm")
// 	}

// 	// Add pagination details
// 	data := models.GitMetricsViewData{
// 		Metrics:           metrics,
// 		Repos:             nil,
// 		Authors:           nil,
// 		RepoNameParameter: repoName,
// 		UserNameParameter: userName,
// 		CurrentPage:       page,
// 		PreviousPage:      page - 1,
// 		NextPage:          page + 1,
// 	}

// 	if err := helpers.RenderTemplateFunc(w, data, "git_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
