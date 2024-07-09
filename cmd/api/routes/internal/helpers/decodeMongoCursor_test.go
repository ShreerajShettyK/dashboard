// package helpers

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// // MockCursor is a mock implementation of mongo.Cursor
// type MockCursor struct {
// 	mock.Mock
// }

// func (m *MockCursor) Next(ctx context.Context) bool {
// 	args := m.Called(ctx)
// 	return args.Bool(0)
// }

// func (m *MockCursor) Decode(val interface{}) error {
// 	args := m.Called(val)
// 	return args.Error(0)
// }

// // Implement other methods of mongo.Cursor interface with empty implementations
// func (m *MockCursor) Close(ctx context.Context) error {
// 	return nil
// }

// func (m *MockCursor) Err() error {
// 	return nil
// }

// func (m *MockCursor) ID() int64 {
// 	return 0
// }

// func (m *MockCursor) All(ctx context.Context, results interface{}) error {
// 	return nil
// }

// func (m *MockCursor) TryNext(ctx context.Context) bool {
// 	return false
// }

// func (m *MockCursor) RemainingBatchLength() int {
// 	return 0
// }

// func (m *MockCursor) Current() bson.Raw {
// 	return nil
// }

// // TestDocument is a sample struct for testing
// type TestDocument struct {
// 	ID   string `bson:"_id"`
// 	Name string `bson:"name"`
// }

// func TestDecodeCursor(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		setupMock      func(*MockCursor)
// 		expectedResult []TestDocument
// 		expectedError  error
// 	}{
// 		{
// 			name: "Successful decoding of multiple documents",
// 			setupMock: func(m *MockCursor) {
// 				m.On("Next", mock.Anything).Return(true).Times(3)
// 				m.On("Next", mock.Anything).Return(false).Once()
// 				m.On("Decode", mock.AnythingOfType("*helpers.TestDocument")).Run(func(args mock.Arguments) {
// 					arg := args.Get(0).(*TestDocument)
// 					*arg = TestDocument{ID: "1", Name: "Test1"}
// 				}).Return(nil).Once()
// 				m.On("Decode", mock.AnythingOfType("*helpers.TestDocument")).Run(func(args mock.Arguments) {
// 					arg := args.Get(0).(*TestDocument)
// 					*arg = TestDocument{ID: "2", Name: "Test2"}
// 				}).Return(nil).Once()
// 				m.On("Decode", mock.AnythingOfType("*helpers.TestDocument")).Run(func(args mock.Arguments) {
// 					arg := args.Get(0).(*TestDocument)
// 					*arg = TestDocument{ID: "3", Name: "Test3"}
// 				}).Return(nil).Once()
// 			},
// 			expectedResult: []TestDocument{
// 				{ID: "1", Name: "Test1"},
// 				{ID: "2", Name: "Test2"},
// 				{ID: "3", Name: "Test3"},
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Empty result set",
// 			setupMock: func(m *MockCursor) {
// 				m.On("Next", mock.Anything).Return(false).Once()
// 			},
// 			expectedResult: []TestDocument{},
// 			expectedError:  nil,
// 		},
// 		{
// 			name: "Decode error",
// 			setupMock: func(m *MockCursor) {
// 				m.On("Next", mock.Anything).Return(true).Once()
// 				m.On("Decode", mock.AnythingOfType("*helpers.TestDocument")).Return(errors.New("decode error")).Once()
// 			},
// 			expectedResult: nil,
// 			expectedError:  errors.New("decode error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockCursor := new(MockCursor)
// 			tt.setupMock(mockCursor)

// 			ctx := context.Background()
// 			results, err := DecodeCursor[TestDocument](ctx, mockCursor)

// 			if tt.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.Equal(t, tt.expectedError.Error(), err.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			assert.Equal(t, tt.expectedResult, results)
// 			mockCursor.AssertExpectations(t)
// 		})
// 	}
// }

// func TestDecodeCursorWithPanic(t *testing.T) {
// 	mockCursor := new(MockCursor)
// 	mockCursor.On("Next", mock.Anything).Return(true).Once()
// 	mockCursor.On("Decode", mock.Anything).Panic("Unexpected panic")

// 	ctx := context.Background()

// 	assert.Panics(t, func() {
// 		_, _ = DecodeCursor[TestDocument](ctx, mockCursor)
// 	})

// 	mockCursor.AssertExpectations(t)
// }

package helpers_test

import (
	"context"
	"dashboard/cmd/api/routes/internal/helpers"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCursor is a mock implementation of helpers.Cursor
type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

// TestDocument is a sample struct for testing
type TestDocument struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func TestDecodeCursor(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockCursor)
		expectedResult []TestDocument
		expectedError  error
	}{
		{
			name: "Successful decoding of multiple documents",
			setupMock: func(m *MockCursor) {
				m.On("Next", mock.Anything).Return(true).Times(3)
				m.On("Next", mock.Anything).Return(false).Once()
				m.On("Decode", mock.AnythingOfType("*helpers_test.TestDocument")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*TestDocument)
					*arg = TestDocument{ID: "1", Name: "Test1"}
				}).Return(nil).Once()
				m.On("Decode", mock.AnythingOfType("*helpers_test.TestDocument")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*TestDocument)
					*arg = TestDocument{ID: "2", Name: "Test2"}
				}).Return(nil).Once()
				m.On("Decode", mock.AnythingOfType("*helpers_test.TestDocument")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*TestDocument)
					*arg = TestDocument{ID: "3", Name: "Test3"}
				}).Return(nil).Once()
			},
			expectedResult: []TestDocument{
				{ID: "1", Name: "Test1"},
				{ID: "2", Name: "Test2"},
				{ID: "3", Name: "Test3"},
			},
			expectedError: nil,
		},
		{
			name: "Empty result set",
			setupMock: func(m *MockCursor) {
				m.On("Next", mock.Anything).Return(false).Once()
			},
			expectedResult: []TestDocument(nil),
			expectedError:  nil,
		},
		{
			name: "Decode error",
			setupMock: func(m *MockCursor) {
				m.On("Next", mock.Anything).Return(true).Once()
				m.On("Decode", mock.AnythingOfType("*helpers_test.TestDocument")).Return(errors.New("decode error")).Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("decode error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCursor := new(MockCursor)
			tt.setupMock(mockCursor)

			ctx := context.Background()
			results, err := helpers.DecodeCursor[TestDocument](ctx, mockCursor)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, results)
			mockCursor.AssertExpectations(t)
		})
	}
}

func TestDecodeCursorWithPanic(t *testing.T) {
	mockCursor := new(MockCursor)
	mockCursor.On("Next", mock.Anything).Return(true).Once()
	mockCursor.On("Decode", mock.Anything).Panic("Unexpected panic")

	ctx := context.Background()

	assert.Panics(t, func() {
		_, _ = helpers.DecodeCursor[TestDocument](ctx, mockCursor)
	})

	mockCursor.AssertExpectations(t)
}
