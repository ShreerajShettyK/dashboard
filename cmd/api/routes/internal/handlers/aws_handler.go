package handlers

import (
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"net/http"
	"time"
)

func FetchAWSMetrics(serviceName string, startDate, endDate time.Time) ([]models.AWSMetric, error) {
	query := `SELECT id, service_name, metric_name, value, timestamp FROM aws_metrics WHERE service_name = ? AND timestamp BETWEEN ? AND ?`
	rows, err := database.DB.Query(query, serviceName, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.AWSMetric
	for rows.Next() {
		var metric models.AWSMetric
		if err := rows.Scan(&metric.ID, &metric.ServiceName, &metric.MetricName, &metric.Value, &metric.Timestamp); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func AWSMetricsHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service_name")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	metrics, err := FetchAWSMetrics(serviceName, startDate, endDate)
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
