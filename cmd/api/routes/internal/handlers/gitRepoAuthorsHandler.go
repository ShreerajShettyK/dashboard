package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func FetchAuthorsByRepo(repoName string) ([]string, error) {
	if GitMetricsCollection == nil {
		return nil, errors.New("GitMetricsCollection is not initialized")
	}

	filter := bson.M{"reponame": repoName}
	cursor, err := GitMetricsCollection.Distinct(context.Background(), "commited_by", filter)
	if err != nil {
		return nil, err
	}

	var authors []string
	for _, author := range cursor {
		if authorStr, ok := author.(string); ok {
			authors = append(authors, authorStr)
		} else {
			return nil, errors.New("type assertion failed for author")
		}
	}
	return authors, nil
}

func GitAuthorsByRepoHandler(w http.ResponseWriter, r *http.Request) {
	repoName := r.URL.Query().Get("repo_name")
	if repoName == "" {
		http.Error(w, "repository name is required", http.StatusBadRequest)
		return
	}

	authors, err := FetchAuthorsByRepo(repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
