package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dashboard/cmd/api/routes/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// jsonEncoder is a type alias for the json.NewEncoder function
// type jsonEncoder func(w io.Writer) *json.Encoder

// var jsonEncoderFunc jsonEncoder = json.NewEncoder

func TestFetchAuthorsByRepo_Success(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedAuthors := []interface{}{"author1", "author2", "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", bson.M{"reponame": "test_repo"}, ([]*options.DistinctOptions)(nil)).Return(expectedAuthors, nil)

	authors, err := handlers.FetchAuthorsByRepo("test_repo")

	assert.NoError(t, err)
	assert.Equal(t, []string{"author1", "author2", "author3"}, authors)
	mockCollection.AssertExpectations(t)
}

func TestFetchAuthorsByRepo_Error(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Distinct", mock.Anything, "commited_by", bson.M{"reponame": "test_repo"}, ([]*options.DistinctOptions)(nil)).Return([]interface{}{}, expectedError)

	authors, err := handlers.FetchAuthorsByRepo("test_repo")

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}

func TestFetchAuthorsByRepo_UninitializedCollection(t *testing.T) {
	handlers.GitMetricsCollection = nil

	authors, err := handlers.FetchAuthorsByRepo("test_repo")

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, "GitMetricsCollection is not initialized", err.Error())
}

func TestFetchAuthorsByRepo_TypeAssertionFailure(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	invalidAuthors := []interface{}{"author1", 123, "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", bson.M{"reponame": "test_repo"}, ([]*options.DistinctOptions)(nil)).Return(invalidAuthors, nil)

	authors, err := handlers.FetchAuthorsByRepo("test_repo")

	assert.Error(t, err)
	assert.Nil(t, authors)
	assert.Equal(t, "type assertion failed for author", err.Error())
	mockCollection.AssertExpectations(t)
}

func TestGitAuthorsByRepoHandler_Success(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedAuthors := []interface{}{"author1", "author2", "author3"}
	mockCollection.On("Distinct", mock.Anything, "commited_by", bson.M{"reponame": "test_repo"}, ([]*options.DistinctOptions)(nil)).Return(expectedAuthors, nil)

	req, err := http.NewRequest("GET", "/git-authors-by-repo?repo_name=test_repo", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitAuthorsByRepoHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var authors []string
	err = json.NewDecoder(rr.Body).Decode(&authors)
	assert.NoError(t, err)
	assert.Equal(t, []string{"author1", "author2", "author3"}, authors)
	mockCollection.AssertExpectations(t)
}

func TestGitAuthorsByRepoHandler_MissingRepoName(t *testing.T) {
	req, err := http.NewRequest("GET", "/git-authors-by-repo", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitAuthorsByRepoHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "repository name is required\n", rr.Body.String())
}

func TestGitAuthorsByRepoHandler_FetchError(t *testing.T) {
	mockCollection := new(MockGitMetricsCollection)
	handlers.GitMetricsCollection = mockCollection

	expectedError := errors.New("database error")
	mockCollection.On("Distinct", mock.Anything, "commited_by", bson.M{"reponame": "test_repo"}, ([]*options.DistinctOptions)(nil)).Return([]interface{}{}, expectedError)

	req, err := http.NewRequest("GET", "/git-authors-by-repo?repo_name=test_repo", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GitAuthorsByRepoHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "database error\n", rr.Body.String())
	mockCollection.AssertExpectations(t)
}

// func TestGitAuthorsByRepoHandler_JSONEncodeError(t *testing.T) {
// 	originalJSONEncoderFunc := jsonEncoderFunc
// 	defer func() { jsonEncoderFunc = originalJSONEncoderFunc }()

// 	// Mock the jsonEncoderFunc to simulate a JSON encoding error
// 	jsonEncoderFunc = func(w io.Writer) *json.Encoder {
// 		return json.NewEncoder(errorWriter{w})
// 	}

// 	mockCollection := new(MockGitMetricsCollection)
// 	handlers.GitMetricsCollection = mockCollection

// 	expectedAuthors := []interface{}{"author1", "author2", "author3"}
// 	mockCollection.On("Distinct", mock.Anything, "commited_by", bson.M{"reponame": "test_repo"}, ([]*options.DistinctOptions)(nil)).Return(expectedAuthors, nil)

// 	req, err := http.NewRequest("GET", "/git-authors-by-repo?repo_name=test_repo", nil)
// 	assert.NoError(t, err)

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlers.GitAuthorsByRepoHandler)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 	assert.Equal(t, "intentional error\n", rr.Body.String())
// 	mockCollection.AssertExpectations(t)
// }

// type errorWriter struct {
// 	io.Writer
// }

// func (ew errorWriter) Write(p []byte) (n int, err error) {
// 	return 0, errors.New("intentional error")
// }
