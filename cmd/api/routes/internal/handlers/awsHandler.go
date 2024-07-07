package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchAWSMetrics(serviceName string, startDate, endDate time.Time) ([]models.AWSMetric, error) {
	filter := bson.M{
		"service_name": serviceName,
		"timestamp": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	opts := options.Find().SetSort(bson.D{{"timestamp", 1}})

	cursor, err := database.AWSMetricsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var metrics []models.AWSMetric
	for cursor.Next(context.Background()) {
		var metric models.AWSMetric
		if err := cursor.Decode(&metric); err != nil {
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

	tmplPath, err := filepath.Abs("templates/git_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, metrics)
}
