package handlers

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GitMetricsCollectionInterface defines the methods we need from the MongoDB collection
type GitMetricsCollectionInterface interface {
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
}

// MockGitMetricsCollection is a mock implementation of GitMetricsCollectionInterface
type MockGitMetricsCollection struct {
	mock.Mock
}

func (m *MockGitMetricsCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	args := m.Called(ctx, fieldName, filter, opts)
	return args.Get(0).([]interface{}), args.Error(1)
}
