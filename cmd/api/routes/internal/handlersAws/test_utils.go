package handlersAws

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GitMetricsCollectionInterface defines the methods we need from the MongoDB collection
type AWSMetricsCollectionInterface interface {
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
}

// MockGitMetricsCollection is a mock implementation of GitMetricsCollectionInterface
type MockAWSMetricsCollection struct {
	mock.Mock
}

func (m *MockAWSMetricsCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	args := m.Called(ctx, fieldName, filter, opts)
	return args.Get(0).([]interface{}), args.Error(1)
}
