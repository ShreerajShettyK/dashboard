// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"net/http"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// func FetchDistinctRepositories() ([]string, error) {
// 	cursor, err := database.GitMetricsCollection.Distinct(context.Background(), "reponame", bson.D{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var repos []string
// 	for _, repo := range cursor {
// 		repos = append(repos, repo.(string))
// 	}
// 	return repos, nil
// }

// // Git Repositories Handler
// func GitReposHandler(w http.ResponseWriter, r *http.Request) {
// 	repos, err := FetchDistinctRepositories()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	data := models.GitMetricsViewData{
// 		Metrics: nil,
// 		Repos:   repos,
// 		Authors: nil,
// 	}

// 	if err := helpers.RenderTemplate(w, data, "git_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"errors"
// 	"net/http"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// func FetchDistinctRepositories() ([]string, error) {
// 	if database.GitMetricsCollection == nil {
// 		return nil, errors.New("GitMetricsCollection is not initialized")
// 	}

// 	cursor, err := database.GitMetricsCollection.Distinct(context.Background(), "reponame", bson.D{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var repos []string
// 	for _, repo := range cursor {
// 		if repoStr, ok := repo.(string); ok {
// 			repos = append(repos, repoStr)
// 		} else {
// 			return nil, errors.New("type assertion failed for repository")
// 		}
// 	}
// 	return repos, nil
// }

// // Git Repositories Handler
// func GitReposHandler(w http.ResponseWriter, r *http.Request) {
// 	repos, err := FetchDistinctRepositories()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	data := models.GitMetricsViewData{
// 		Metrics: nil,
// 		Repos:   repos,
// 		Authors: nil,
// 	}

// 	if err := helpers.RenderTemplateFunc(w, data, "git_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// testing prev one works

package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func FetchDistinctRepositories() ([]string, error) {
	if GitMetricsCollection == nil {
		return nil, errors.New("GitMetricsCollection is not initialized")
	}

	cursor, err := GitMetricsCollection.Distinct(context.Background(), "reponame", bson.D{})
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, repo := range cursor {
		if repoStr, ok := repo.(string); ok {
			repos = append(repos, repoStr)
		} else {
			return nil, errors.New("type assertion failed for repository")
		}
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

	if err := helpers.RenderTemplateFunc(w, data, "git_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
