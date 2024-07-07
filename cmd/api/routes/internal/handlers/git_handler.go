package handlers

import (
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"net/http"
)

func FetchGitMetrics(userName, repoName string) ([]models.GitMetric, error) {
	query := `SELECT id, repo_name, user_name, commit_id, commit_date, lines_added, lines_deleted, files_added, files_deleted FROM git_metrics WHERE user_name = ? AND repo_name = ?`
	rows, err := database.DB.Query(query, userName, repoName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.GitMetric
	for rows.Next() {
		var metric models.GitMetric
		if err := rows.Scan(&metric.ID, &metric.RepoName, &metric.UserName, &metric.CommitID, &metric.CommitDate, &metric.LinesAdded, &metric.LinesDeleted, &metric.FilesAdded, &metric.FilesDeleted); err != nil {
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

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, metrics)
}
