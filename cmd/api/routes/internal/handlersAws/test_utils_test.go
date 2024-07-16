package handlersAws

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMockAWSMetricsCollection_Distinct(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockAWSMetricsCollection)
		expectedResult []interface{}
		expectedError  error
	}{
		{
			name: "Successful Distinct Call",
			setupMock: func(m *MockAWSMetricsCollection) {
				m.On("Distinct", mock.Anything, "fieldName", mock.Anything, mock.Anything).
					Return([]interface{}{"value1", "value2"}, nil)
			},
			expectedResult: []interface{}{"value1", "value2"},
			expectedError:  nil,
		},
		{
			name: "Distinct Call with Error",
			setupMock: func(m *MockAWSMetricsCollection) {
				m.On("Distinct", mock.Anything, "fieldName", mock.Anything, mock.Anything).
					Return([]interface{}{}, errors.New("database error"))
			},
			expectedResult: []interface{}{},
			expectedError:  errors.New("database error"),
		},
		{
			name: "Distinct Call with Empty Result",
			setupMock: func(m *MockAWSMetricsCollection) {
				m.On("Distinct", mock.Anything, "fieldName", mock.Anything, mock.Anything).
					Return([]interface{}{}, nil)
			},
			expectedResult: []interface{}{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := new(MockAWSMetricsCollection)
			tt.setupMock(mockCollection)

			ctx := context.Background()
			filter := map[string]string{"key": "value"}
			opts := options.Distinct()

			result, err := mockCollection.Distinct(ctx, "fieldName", filter, opts)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockCollection.AssertExpectations(t)
		})
	}
}

func TestAWSMetricsCollectionInterface(t *testing.T) {
	var _ AWSMetricsCollectionInterface = (*MockAWSMetricsCollection)(nil)
}
