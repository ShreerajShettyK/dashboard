package handlersAws_test

import (
	"dashboard/cmd/api/routes/internal/handlersAws"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwsHomeHandler_Success(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	// Mock the RenderTemplateFunc
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		assert.Equal(t, "aws_dashboard.html", templateName)
		assert.IsType(t, models.AwsMetricsViewData{}, data)
		viewData := data.(models.AwsMetricsViewData)
		assert.Equal(t, []string{"ec2", "rds", "elb"}, viewData.Services)
		assert.Nil(t, viewData.Metrics)
		return nil
	}

	req, err := http.NewRequest("GET", "/aws-home", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlersAws.AwsHomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAwsHomeHandler_RenderError(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	// Mock the RenderTemplateFunc to return an error
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		return errors.New("render error")
	}

	req, err := http.NewRequest("GET", "/aws-home", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlersAws.AwsHomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "render error\n", rr.Body.String())
}
