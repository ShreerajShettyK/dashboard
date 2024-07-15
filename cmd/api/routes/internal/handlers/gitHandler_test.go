package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dashboard/cmd/api/routes/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewMockGitMetricsCollection struct {
	mock.Mock
}

type MockCursor struct {
	mock.Mock
}

// Find implements handlers.NewGitMetricsCollectionInterface
func (m *NewMockGitMetricsCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// func TestFetchGitMetrics_Success(t *testing.T) {
// 	mockCollection := new(NewMockGitMetricsCollection)
// 	handlers.NewGitMetricsCollection = mockCollection

// 	expectedMetrics := []models.GitMetric{
// 		{CommitDate: time.Now()},
// 		{CommitDate: time.Now()},
// 	}

// 	mockCursor := &mongo.Cursor{}

// 	mockCollection.On("Find", mock.Anything, bson.M{"commited_by": "test_user", "reponame": "test_repo"}, mock.Anything).Return(mockCursor, nil)

// 	metrics, err := handlers.FetchGitMetrics("test_user", "test_repo", 10, 0)

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedMetrics, metrics)
// 	mockCollection.AssertExpectations(t)
// }

// ... (rest of the test cases)
func TestFetchGitMetrics_Error(t *testing.T) {
	mockCollection := new(NewMockGitMetricsCollection)
	handlers.NewGitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Find", mock.Anything, bson.M{"commited_by": "test_user", "reponame": "test_repo"}, mock.Anything).Return(nil, expectedError)

	metrics, err := handlers.FetchGitMetrics("test_user", "test_repo", 10, 0)

	assert.Error(t, err)
	assert.Nil(t, metrics)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}

// func TestGitMetricsHandler_Success(t *testing.T) {
// 	mockCollection := new(NewMockGitMetricsCollection)
// 	handlers.NewGitMetricsCollection = mockCollection

// 	expectedMetrics := []models.GitMetric{
// 		{CommitDate: time.Now()},
// 		{CommitDate: time.Now()},
// 	}

// 	mockCursor := &mongo.Cursor{}

// 	mockCollection.On("Find", mock.Anything, bson.M{"commited_by": "test_user", "reponame": "test_repo"}, mock.Anything).Return(mockCursor, nil)

// 	// Mock the cursor.All method to return the expectedMetrics
// 	mockCollection.On("All", mock.Anything, mock.AnythingOfType("*[]models.GitMetric")).Run(func(args mock.Arguments) {
// 		arg := args.Get(1).(*[]models.GitMetric)
// 		*arg = expectedMetrics
// 	}).Return(nil)

// 	req, err := http.NewRequest("GET", "/git-metrics?user_name=test_user&repo_name=test_repo&page=1", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlers.GitMetricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	// Verify that the expected metrics were used
// 	var result models.GitMetricsViewData
// 	err = json.NewDecoder(rr.Body).Decode(&result)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedMetrics, result.Metrics)

// 	mockCollection.AssertExpectations(t)
// }

func TestGitMetricsHandler_FetchError(t *testing.T) {
	mockCollection := new(NewMockGitMetricsCollection)
	handlers.NewGitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Find", mock.Anything, bson.M{"commited_by": "test_user", "reponame": "test_repo"}, mock.Anything).Return(nil, expectedError)

	req, err := http.NewRequest("GET", "/git-metrics?user_name=test_user&repo_name=test_repo&page=1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitMetricsHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "database error\n", rr.Body.String())
	mockCollection.AssertExpectations(t)
}

// func TestGitMetricsHandler_RenderTemplateError(t *testing.T) {
// 	originalRenderTemplateFunc := helpers.RenderTemplateFunc
// 	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

// 	mockCollection := new(NewMockGitMetricsCollection)
// 	handlers.NewGitMetricsCollection = mockCollection

// 	expectedMetrics := []models.GitMetric{
// 		{CommitDate: time.Now()},
// 		{CommitDate: time.Now()},
// 	}

// 	mockCursor := new(MockCursor)
// 	mockCursor.On("Close", mock.Anything).Return(nil)
// 	mockCursor.On("Next", mock.Anything).Return(true).Once()
// 	mockCursor.On("Next", mock.Anything).Return(false).Once()
// 	mockCursor.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
// 		arg := args.Get(0).(*models.GitMetric)
// 		*arg = expectedMetrics[0]
// 		expectedMetrics = expectedMetrics[1:]
// 	})

// 	mockCollection.On("Find", mock.Anything, bson.M{"commited_by": "test_user", "reponame": "test_repo"}, mock.Anything).Return(mockCursor, nil)

// 	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
// 		return errors.New("render error")
// 	}

// 	req, err := http.NewRequest("GET", "/git-metrics?user_name=test_user&repo_name=test_repo&page=1", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlers.GitMetricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 	assert.Equal(t, "template error\n", rr.Body.String())
// 	mockCollection.AssertExpectations(t)
// }

// func TestGitMetricsHandler_InvalidPage(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/git-metrics?user_name=test_user&repo_name=test_repo&page=invalid", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlers.GitMetricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestGitMetricsHandler_NoQueryParameters(t *testing.T) {
// 	mockCollection := new(NewMockGitMetricsCollection)
// 	handlers.NewGitMetricsCollection = mockCollection

// 	expectedMetrics := []models.GitMetric{
// 		{CommitDate: time.Now()},
// 		{CommitDate: time.Now()},
// 	}

// 	mockCursor := new(MockCursor)
// 	mockCursor.On("Close", mock.Anything).Return(nil)
// 	mockCursor.On("Next", mock.Anything).Return(true).Once()
// 	mockCursor.On("Next", mock.Anything).Return(false).Once()
// 	mockCursor.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
// 		arg := args.Get(0).(*models.GitMetric)
// 		*arg = expectedMetrics[0]
// 		expectedMetrics = expectedMetrics[1:]
// 	})

// 	mockCollection.On("Find", mock.Anything, bson.M{}, mock.Anything).Return(mockCursor, nil)

// 	req, err := http.NewRequest("GET", "/git-metrics?page=1", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlers.GitMetricsHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	mockCollection.AssertExpectations(t)
// }
