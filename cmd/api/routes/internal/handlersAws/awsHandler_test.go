package handlersAws_test

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/handlersAws"
// 	"dashboard/cmd/api/routes/internal/models"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type MockAWSMetricsCollection struct {
// 	mock.Mock
// }

// func (m *MockAWSMetricsCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
// 	args := m.Called(ctx, fieldName, filter, opts)
// 	return args.Get(0).([]interface{}), args.Error(1)
// }

// func TestFetchAWSMetrics_Success(t *testing.T) {
// 	mockCollection := new(MockAWSMetricsCollection)
// 	handlersAws.AWSMetricsCollection = mockCollection

// 	expectedMetrics := []models.AWSMetric{{ServiceName: "ec2", Date: time.Now()}}
// 	cursor := new(mongo.Cursor) // Mock cursor
// 	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(cursor, nil)

// 	metrics, err := handlersAws.FetchAWSMetrics("ec2", time.Now(), time.Now())

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedMetrics, metrics)
// 	mockCollection.AssertExpectations(t)
// }

// func TestFetchAWSMetrics_Error(t *testing.T) {
// 	mockCollection := new(MockAWSMetricsCollection)
// 	handlersAws.AWSMetricsCollection = mockCollection

// 	expectedError := errors.New("database error")
// 	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

// 	metrics, err := handlersAws.FetchAWSMetrics("ec2", time.Now(), time.Now())

// 	assert.Error(t, err)
// 	assert.Nil(t, metrics)
// 	assert.Equal(t, expectedError, err)
// 	mockCollection.AssertExpectations(t)
// }

// // func TestAWSMetricsHandler_Success(t *testing.T) {
// // 	originalRenderTemplateFunc := helpers.RenderTemplateFunc
// // 	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

// // 	mockCollection := new(MockAWSMetricsCollection)
// // 	handlersAws.AWSMetricsCollection = mockCollection

// // 	expectedMetrics := []models.AWSMetric{{ServiceName: "ec2", Date: time.Now()}}
// // 	cursor := new(mongo.Cursor) // Mock cursor
// // 	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(cursor, nil)
// // 	helpers.DecodeCursorFunc = func(ctx context.Context, cursor *mongo.Cursor) ([]models.AWSMetric, error) {
// // 		return expectedMetrics, nil
// // 	}

// // 	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
// // 		assert.Equal(t, "aws_dashboard.html", templateName)
// // 		assert.IsType(t, models.AwsMetricsViewData{}, data)
// // 		viewData := data.(models.AwsMetricsViewData)
// // 		assert.Equal(t, expectedMetrics, viewData.Metrics)
// // 		assert.Equal(t, []string{"ec2", "rds", "elb"}, viewData.Services)
// // 		return nil
// // 	}

// // 	req, err := http.NewRequest("GET", "/aws-metrics?service_name=ec2&date_range=2023-01-01+-+2023-01-31", nil)
// // 	assert.NoError(t, err)

// // 	rr := httptest.NewRecorder()
// // 	handler := http.HandlerFunc(handlersAws.AWSMetricsHandler)

// // 	handler.ServeHTTP(rr, req)

// // 	assert.Equal(t, http.StatusOK, rr.Code)
// // }

// func TestAWSMetricsHandler_FetchError(t *testing.T) {
// 	mockCollection := new(MockAWSMetricsCollection)
// 	handlersAws.AWSMetricsCollection = mockCollection

// 	expectedError := errors.New("database error")
// 	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

// 	req, err := http.NewRequest("GET", "/aws-metrics?service_name=ec2&date_range=2023-01-01+-+2023-01-31", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlersAws.AWSMetricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 	assert.Equal(t, "database error\n", rr.Body.String())
// }

// // func TestAWSMetricsHandler_RenderError(t *testing.T) {
// // 	originalRenderTemplateFunc := helpers.RenderTemplateFunc
// // 	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

// // 	mockCollection := new(MockAWSMetricsCollection)
// // 	handlersAws.AWSMetricsCollection = mockCollection

// // 	expectedMetrics := []models.AWSMetric{{ServiceName: "ec2", Date: time.Now()}}
// // 	cursor := new(mongo.Cursor) // Mock cursor
// // 	mockCollection.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(cursor, nil)
// // 	helpers.DecodeCursorFunc = func(ctx context.Context, cursor *mongo.Cursor) ([]models.AWSMetric, error) {
// // 		return expectedMetrics, nil
// // 	}

// // 	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
// // 		return errors.New("render error")
// // 	}

// // 	req, err := http.NewRequest("GET", "/aws-metrics?service_name=ec2&date_range=2023-01-01+-+2023-01-31", nil)
// // 	assert.NoError(t, err)

// // 	rr := httptest.NewRecorder()
// // 	handler := http.HandlerFunc(handlersAws.AWSMetricsHandler)

// // 	handler.ServeHTTP(rr, req)

// // 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// // 	assert.Equal(t, "render error\n", rr.Body.String())
// // }

// func TestAWSMetricsHandler_InvalidDateRange(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/aws-metrics?service_name=ec2&date_range=invalid-date", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlersAws.AWSMetricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// }
