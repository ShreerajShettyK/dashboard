package handlers

import (
	"dashboard/cmd/api/routes/internal/models"
	"html/template"
	"net/http"
)

func AwsHomeHandler(w http.ResponseWriter, r *http.Request) {
	services := []string{"ec2", "rds", "elb"}

	data := models.AwsMetricsViewData{
		Metrics:  nil,
		Services: services,
	}

	// Use the relative path
	tmplPath := "internal/templates/aws_dashboard.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with repos data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
