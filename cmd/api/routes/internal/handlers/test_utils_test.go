package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMockGitMetricsCollection_Distinct(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockGitMetricsCollection)
		expectedResult []interface{}
		expectedError  error
	}{
		{
			name: "Successful Distinct Call",
			setupMock: func(m *MockGitMetricsCollection) {
				m.On("Distinct", mock.Anything, "fieldName", mock.Anything, mock.Anything).
					Return([]interface{}{"value1", "value2"}, nil)
			},
			expectedResult: []interface{}{"value1", "value2"},
			expectedError:  nil,
		},
		{
			name: "Distinct Call with Error",
			setupMock: func(m *MockGitMetricsCollection) {
				m.On("Distinct", mock.Anything, "fieldName", mock.Anything, mock.Anything).
					Return([]interface{}{}, errors.New("database error"))
			},
			expectedResult: []interface{}{},
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := new(MockGitMetricsCollection)
			tt.setupMock(mockCollection)

			ctx := context.Background()
			filter := map[string]string{"key": "value"}
			opts := options.Distinct()

			result, err := mockCollection.Distinct(ctx, "fieldName", filter, opts)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
			mockCollection.AssertExpectations(t)
		})
	}
}
