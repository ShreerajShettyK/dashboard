package handlers_test

import (
	"context"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dashboard/cmd/api/routes/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockGitMetricsCollection struct {
	mock.Mock
}

func (m *MockGitMetricsCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	args := m.Called(ctx, fieldName, filter, opts)
	return args.Get(0).([]interface{}), args.Error(1)
}

func TestFetchDistinctAuthors_Success(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedAuthors := []interface{}{"author1", "author2", "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", mock.Anything, mock.Anything).Return(expectedAuthors, nil)

	authors, err := handlers.FetchDistinctAuthors()

	assert.NoError(t, err)
	assert.Equal(t, []string{"author1", "author2", "author3"}, authors)
	mockCollection.AssertExpectations(t)
}

func TestFetchDistinctAuthors_Error(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Distinct", mock.Anything, "commited_by", mock.Anything, mock.Anything).Return([]interface{}{}, expectedError)

	authors, err := handlers.FetchDistinctAuthors()

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}

func TestFetchDistinctAuthors_UninitializedCollection(t *testing.T) {
	handlers.GitMetricsCollection = nil

	authors, err := handlers.FetchDistinctAuthors()

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, "GitMetricsCollection is not initialized", err.Error())
}

func TestFetchDistinctAuthors_TypeAssertionFailure(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	invalidAuthors := []interface{}{"author1", 123, "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", mock.Anything, mock.Anything).Return(invalidAuthors, nil)

	authors, err := handlers.FetchDistinctAuthors()

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, "type assertion failed for author", err.Error())
	mockCollection.AssertExpectations(t)
}

func TestGitAuthorsHandler_FetchError(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Distinct", mock.Anything, "commited_by", mock.Anything, mock.Anything).Return([]interface{}{}, expectedError)

	req, err := http.NewRequest("GET", "/git-authors", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitAuthorsHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockCollection.AssertExpectations(t)
}

func TestGitAuthorsHandler_Success(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedAuthors := []interface{}{"author1", "author2", "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", mock.Anything, mock.Anything).Return(expectedAuthors, nil)

	// Mock the RenderTemplateFunc
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		assert.Equal(t, "git_dashboard.html", templateName)
		assert.IsType(t, models.GitMetricsViewData{}, data)
		viewData := data.(models.GitMetricsViewData)
		assert.Equal(t, []string{"author1", "author2", "author3"}, viewData.Authors)
		return nil
	}

	req, err := http.NewRequest("GET", "/git-authors", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitAuthorsHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCollection.AssertExpectations(t)
}

func TestGitAuthorsHandler_RenderError(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedAuthors := []interface{}{"author1", "author2", "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", mock.Anything, mock.Anything).Return(expectedAuthors, nil)

	// Mock the RenderTemplateFunc to return an error
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		return errors.New("render error")
	}

	req, err := http.NewRequest("GET", "/git-authors", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitAuthorsHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockCollection.AssertExpectations(t)
}
