package handlers_test

import (
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dashboard/cmd/api/routes/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchDistinctRepositories_Success(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedRepos := []interface{}{"repo1", "repo2", "repo3"}
	mockCollection.On("Distinct", mock.Anything, "reponame", mock.Anything, mock.Anything).Return(expectedRepos, nil)

	repos, err := handlers.FetchDistinctRepositories()

	assert.NoError(t, err)
	assert.Equal(t, []string{"repo1", "repo2", "repo3"}, repos)
	mockCollection.AssertExpectations(t)
}

func TestFetchDistinctRepositories_Error(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Distinct", mock.Anything, "reponame", mock.Anything, mock.Anything).Return([]interface{}{}, expectedError)

	repos, err := handlers.FetchDistinctRepositories()

	assert.Error(t, err)
	assert.Nil(t, repos)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}

// func TestFetchDistinctRepositories_UninitializedCollection(t *testing.T) {
// 	database.GitMetricsCollection = nil

// 	repos, err := handlers.FetchDistinctRepositories()

// 	assert.Error(t, err)
// 	assert.Nil(t, repos)
// 	assert.Equal(t, "GitMetricsCollection is not initialized", err.Error())
// }

func TestFetchDistinctRepositories_UninitializedCollection(t *testing.T) {
	handlers.GitMetricsCollection = nil

	authors, err := handlers.FetchDistinctRepositories()

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, "GitMetricsCollection is not initialized", err.Error())
}

func TestFetchDistinctRepositories_TypeAssertionFailure(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	invalidRepos := []interface{}{"repo1", 123, "repo3"}
	mockCollection.On("Distinct", mock.Anything, "reponame", mock.Anything, mock.Anything).Return(invalidRepos, nil)

	repos, err := handlers.FetchDistinctRepositories()

	assert.Error(t, err)
	assert.Nil(t, repos)
	assert.Equal(t, "type assertion failed for repository", err.Error())
	mockCollection.AssertExpectations(t)
}

func TestGitReposHandler_FetchError(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Distinct", mock.Anything, "reponame", mock.Anything, mock.Anything).Return([]interface{}{}, expectedError)

	req, err := http.NewRequest("GET", "/git-repos", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitReposHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockCollection.AssertExpectations(t)
}

func TestGitReposHandler_Success(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedRepos := []interface{}{"repo1", "repo2", "repo3"}
	mockCollection.On("Distinct", mock.Anything, "reponame", mock.Anything, mock.Anything).Return(expectedRepos, nil)

	// Mock the RenderTemplateFunc
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		assert.Equal(t, "git_dashboard.html", templateName)
		assert.IsType(t, models.GitMetricsViewData{}, data)
		viewData := data.(models.GitMetricsViewData)
		assert.Equal(t, []string{"repo1", "repo2", "repo3"}, viewData.Repos)
		return nil
	}

	req, err := http.NewRequest("GET", "/git-repos", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitReposHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockCollection.AssertExpectations(t)
}

func TestGitReposHandler_RenderError(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedRepos := []interface{}{"repo1", "repo2", "repo3"}
	mockCollection.On("Distinct", mock.Anything, "reponame", mock.Anything, mock.Anything).Return(expectedRepos, nil)

	// Mock the RenderTemplateFunc to return an error
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		return errors.New("render error")
	}

	req, err := http.NewRequest("GET", "/git-repos", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitReposHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockCollection.AssertExpectations(t)
}
