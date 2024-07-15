package routes

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	args := m.Called(path, f)
	return args.Get(0).(*mux.Route)
}

func (m *MockRouter) Methods(methods ...string) *mux.Route {
	args := m.Called(methods)
	return args.Get(0).(*mux.Route)
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "Router should not be nil")
	assert.IsType(t, &mux.Router{}, router, "Router should be of type *mux.Router")

	// Test that all routes are registered
	routes := []struct {
		path   string
		method string
	}{
		{"/aws_metrics/home", "GET"},
		{"/aws_metrics/home/resources", "GET"},
		{"/git_metrics/home", "GET"},
		{"/git_metrics/home/commits", "GET"},
		{"/git_metrics/home/repos", "GET"},
		{"/git_metrics/home/authors", "GET"},
		{"/git_metrics/repoAuthors", "GET"},
	}

	for _, route := range routes {
		req, _ := http.NewRequest(route.method, route.path, nil)
		match := &mux.RouteMatch{}
		matched := router.Match(req, match)
		assert.True(t, matched, "Route %s should be registered", route.path)
	}

	// Negative test case: non-existent route
	req, _ := http.NewRequest("GET", "/non_existent_route", nil)
	match := &mux.RouteMatch{}
	matched := router.Match(req, match)
	assert.False(t, matched, "Non-existent route should not match")

	// Negative test case: incorrect method
	req, _ = http.NewRequest("POST", "/aws_metrics/home", nil)
	match = &mux.RouteMatch{}
	matched = router.Match(req, match)
	assert.False(t, matched, "Route with incorrect method should not match")
}
