package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dashboard/cmd/api/routes/internal/handlers"
	"dashboard/cmd/api/routes/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGitHomeHandler_Success(t *testing.T) {
	// Mock the FetchDistinctAuthors and FetchDistinctRepositories functions
	originalFetchDistinctAuthors := handlers.FetchDistinctAuthorsMock
	originalFetchDistinctRepositories := handlers.FetchDistinctRepositoriesMock
	originalRenderTemplate := handlers.RenderTemplate

	defer func() {
		handlers.FetchDistinctAuthorsMock = originalFetchDistinctAuthors
		handlers.FetchDistinctRepositoriesMock = originalFetchDistinctRepositories
		handlers.RenderTemplate = originalRenderTemplate
	}()

	handlers.FetchDistinctAuthorsMock = func() ([]string, error) {
		return []string{"author1", "author2", "author3"}, nil
	}

	handlers.FetchDistinctRepositoriesMock = func() ([]string, error) {
		return []string{"repo1", "repo2", "repo3"}, nil
	}

	handlers.RenderTemplate = func(w http.ResponseWriter, data interface{}, templateName string) error {
		assert.Equal(t, "git_dashboard.html", templateName)
		assert.IsType(t, models.GitMetricsViewData{}, data)
		viewData := data.(models.GitMetricsViewData)
		assert.Equal(t, []string{"author1", "author2", "author3"}, viewData.Authors)
		assert.Equal(t, []string{"repo1", "repo2", "repo3"}, viewData.Repos)
		return nil
	}

	req, err := http.NewRequest("GET", "/git-home", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitHomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGitHomeHandler_FetchAuthorsError(t *testing.T) {
	// Mock the FetchDistinctAuthors function to return an error
	originalFetchDistinctAuthors := handlers.FetchDistinctAuthorsMock
	originalFetchDistinctRepositories := handlers.FetchDistinctRepositoriesMock

	defer func() {
		handlers.FetchDistinctAuthorsMock = originalFetchDistinctAuthors
		handlers.FetchDistinctRepositoriesMock = originalFetchDistinctRepositories
	}()

	handlers.FetchDistinctAuthorsMock = func() ([]string, error) {
		return nil, errors.New("authors fetch error")
	}

	req, err := http.NewRequest("GET", "/git-home", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitHomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "authors fetch error\n", rr.Body.String())
}

func TestGitHomeHandler_FetchRepositoriesError(t *testing.T) {
	// Mock the FetchDistinctAuthors and FetchDistinctRepositories functions
	originalFetchDistinctAuthors := handlers.FetchDistinctAuthorsMock
	originalFetchDistinctRepositories := handlers.FetchDistinctRepositoriesMock
	originalRenderTemplate := handlers.RenderTemplate

	defer func() {
		handlers.FetchDistinctAuthorsMock = originalFetchDistinctAuthors
		handlers.FetchDistinctRepositoriesMock = originalFetchDistinctRepositories
		handlers.RenderTemplate = originalRenderTemplate
	}()

	handlers.FetchDistinctAuthorsMock = func() ([]string, error) {
		return []string{"author1", "author2", "author3"}, nil
	}

	handlers.FetchDistinctRepositoriesMock = func() ([]string, error) {
		return nil, errors.New("repositories fetch error")
	}

	handlers.RenderTemplate = func(w http.ResponseWriter, data interface{}, templateName string) error {
		// Ensure this doesn't get called
		t.Errorf("RenderTemplate should not be called in this test case")
		return nil
	}

	req, err := http.NewRequest("GET", "/git-home", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitHomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "repositories fetch error\n", rr.Body.String())
}

func TestGitHomeHandler_RenderTemplateError(t *testing.T) {
	// Mock the FetchDistinctAuthors and FetchDistinctRepositories functions and RenderTemplate
	originalFetchDistinctAuthors := handlers.FetchDistinctAuthorsMock
	originalFetchDistinctRepositories := handlers.FetchDistinctRepositoriesMock
	originalRenderTemplate := handlers.RenderTemplate

	defer func() {
		handlers.FetchDistinctAuthorsMock = originalFetchDistinctAuthors
		handlers.FetchDistinctRepositoriesMock = originalFetchDistinctRepositories
		handlers.RenderTemplate = originalRenderTemplate
	}()

	handlers.FetchDistinctAuthorsMock = func() ([]string, error) {
		return []string{"author1", "author2", "author3"}, nil
	}

	handlers.FetchDistinctRepositoriesMock = func() ([]string, error) {
		return []string{"repo1", "repo2", "repo3"}, nil
	}

	handlers.RenderTemplate = func(w http.ResponseWriter, data interface{}, templateName string) error {
		return errors.New("render template error")
	}

	req, err := http.NewRequest("GET", "/git-home", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitHomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "render template error\n", rr.Body.String())
}
