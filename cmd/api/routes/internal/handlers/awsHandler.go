// package handlers

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"dashboard/cmd/api/routes/internal/helpers"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"net/http"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func FetchAWSMetrics(serviceName string, startDate, endDate time.Time) ([]models.AWSMetric, error) {
// 	filter := bson.M{
// 		"service_name": serviceName,
// 		"date": bson.M{
// 			"$gte": startDate,
// 			"$lte": endDate,
// 		},
// 	}
// 	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}})

// 	cursor, err := database.AWSMetricsCollection.Find(context.Background(), filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	return helpers.DecodeCursor[models.AWSMetric](context.Background(), cursor)
// }

// func AWSMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	serviceName := r.URL.Query().Get("service_name")
// 	startDateStr := r.URL.Query().Get("start_date")
// 	endDateStr := r.URL.Query().Get("end_date")

// 	startDate, _ := time.Parse("2006-01-02", startDateStr)
// 	endDate, _ := time.Parse("2006-01-02", endDateStr)

// 	metrics, err := FetchAWSMetrics(serviceName, startDate, endDate)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	services := []string{"ec2", "rds", "elb"}

// 	data := models.AwsMetricsViewData{
// 		Metrics:     metrics,
// 		Services:    services,
// 		ServiceName: serviceName,
// 		StartDate:   startDateStr,
// 		EndDate:     endDateStr,
// 	}

// 	if err := helpers.RenderTemplate(w, data, "aws_dashboard.html"); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

package handlers

import (
	"context"
	"dashboard/cmd/api/routes/internal/database"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchAWSMetrics(serviceName string, startDate, endDate time.Time) ([]models.AWSMetric, error) {
	filter := bson.M{
		"service_name": serviceName,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}})

	cursor, err := database.AWSMetricsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	return helpers.DecodeCursor[models.AWSMetric](context.Background(), cursor)
}

func AWSMetricsHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service_name")
	dateRangeStr := r.URL.Query().Get("date_range")

	var startDate, endDate time.Time
	if dateRangeStr != "" {
		dates := strings.Split(dateRangeStr, " - ")
		startDate, _ = time.Parse("2006-01-02", dates[0])
		endDate, _ = time.Parse("2006-01-02", dates[1])
	} else {
		// Default to a date range of the past month if not specified
		endDate = time.Now()
		startDate = endDate.AddDate(0, -1, 0)
	}

	metrics, err := FetchAWSMetrics(serviceName, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	services := []string{"ec2", "rds", "elb"}

	data := models.AwsMetricsViewData{
		Metrics:     metrics,
		Services:    services,
		ServiceName: serviceName,
		StartDate:   startDate.Format("2006-01-02"),
		EndDate:     endDate.Format("2006-01-02"),
	}

	if err := helpers.RenderTemplate(w, data, "aws_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
