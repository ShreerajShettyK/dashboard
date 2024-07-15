package handlersAws

import (
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"net/http"
)

func AwsHomeHandler(w http.ResponseWriter, r *http.Request) {
	services := []string{"ec2", "rds", "elb"}

	data := models.AwsMetricsViewData{
		Metrics:  nil,
		Services: services,
	}

	if err := helpers.RenderTemplate(w, data, "aws_dashboard.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
