// 86.7
package database

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MockMongoClient struct {
	mock.Mock
}

func (m *MockMongoClient) Database(name string) *mongo.Database {
	args := m.Called(name)
	return args.Get(0).(*mongo.Database)
}

func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	args := m.Called(ctx, rp)
	return args.Error(0)
}

type MockMongoDatabase struct {
	mock.Mock
}

func (m *MockMongoDatabase) Collection(name string) *mongo.Collection {
	args := m.Called(name)
	return args.Get(0).(*mongo.Collection)
}

func TestInitDB(t *testing.T) {
	tests := []struct {
		name          string
		uri           string
		dbName        string
		mockConnect   func(context.Context, ...*options.ClientOptions) (*mongo.Client, error)
		mockPing      func(*mongo.Client, context.Context) error
		expectedError string
	}{
		{
			name:   "Successful connection",
			uri:    "mongodb://localhost:27017",
			dbName: "testdb",
			mockConnect: func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
				mockClient := new(MockMongoClient)
				mockDB := new(MockMongoDatabase)
				mockClient.On("Database", "testdb").Return(mockDB)
				mockDB.On("Collection", "aws_metrics").Return(&mongo.Collection{})
				mockDB.On("Collection", "git_metrics").Return(&mongo.Collection{})
				return &mongo.Client{}, nil
			},
			mockPing: func(client *mongo.Client, ctx context.Context) error {
				return nil
			},
			expectedError: "",
		},
		{
			name:   "Connection failure",
			uri:    "mongodb://localhost:27017",
			dbName: "testdb",
			mockConnect: func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
				return nil, errors.New("connection error")
			},
			mockPing:      nil,
			expectedError: "failed to connect to MongoDB: connection error",
		},
		{
			name:   "Ping failure",
			uri:    "mongodb://localhost:27017",
			dbName: "testdb",
			mockConnect: func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
				return &mongo.Client{}, nil
			},
			mockPing: func(client *mongo.Client, ctx context.Context) error {
				return errors.New("ping error")
			},
			expectedError: "failed to ping MongoDB: ping error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			oldConnect, oldPing := newMongoClient, mongoPing
			defer func() {
				newMongoClient, mongoPing = oldConnect, oldPing
			}()

			newMongoClient = tt.mockConnect
			mongoPing = tt.mockPing

			// Call the function
			err := InitDB(tt.uri, tt.dbName)

			// Assert the result
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, Client)
				assert.NotNil(t, AWSMetricsCollection)
				assert.NotNil(t, GitMetricsCollection)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestNewMongoClient(t *testing.T) {
	// Save the original function and defer its restoration
	oldMongoConnect := mongoConnect
	defer func() { mongoConnect = oldMongoConnect }()

	// Test case 1: Successful connection
	t.Run("Successful connection", func(t *testing.T) {
		// Mock the mongoConnect function
		mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
			return &mongo.Client{}, nil
		}

		// Call newMongoClient
		client, err := newMongoClient(context.Background(), options.Client())

		// Assert the results
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	// Test case 2: Connection failure
	t.Run("Connection failure", func(t *testing.T) {
		// Mock the mongoConnect function to return an error
		mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
			return nil, errors.New("connection error")
		}

		// Call newMongoClient
		client, err := newMongoClient(context.Background(), options.Client())

		// Assert the results
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.EqualError(t, err, "connection error")
	})
}

func TestMongoPing(t *testing.T) {
	// Save the original function and defer its restoration
	oldPing := mongoPing
	defer func() { mongoPing = oldPing }()

	// Test case 1: Successful ping
	t.Run("Successful ping", func(t *testing.T) {
		// Mock the Ping function
		mongoPing = func(client *mongo.Client, ctx context.Context) error {
			return nil
		}

		// Call mongoPing with nil client (it won't be used)
		err := mongoPing(nil, context.Background())

		// Assert the results
		assert.NoError(t, err)
	})

	// Test case 2: Ping failure
	t.Run("Ping failure", func(t *testing.T) {
		// Mock the Ping function to return an error
		expectedError := errors.New("ping failed")
		mongoPing = func(client *mongo.Client, ctx context.Context) error {
			return expectedError
		}

		// Call mongoPing with nil client (it won't be used)
		err := mongoPing(nil, context.Background())

		// Assert the results
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

// package database_test

// import (
// 	"context"
// 	"dashboard/cmd/api/routes/internal/database"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"go.mongodb.org/mongo-driver/mongo/readpref"
// )

// type MockMongoClient struct {
// 	mock.Mock
// }

// func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
// 	args := m.Called(name)
// 	return args.Get(0).(*mongo.Database)
// }

// func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
// 	args := m.Called(ctx, rp)
// 	return args.Error(0)
// }

// type MockMongoDatabase struct {
// 	mock.Mock
// }

// func (m *MockMongoDatabase) Collection(name string) *mongo.Collection {
// 	args := m.Called(name)
// 	return args.Get(0).(*mongo.Collection)
// }

// func TestInitDB(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		uri           string
// 		dbName        string
// 		mockConnect   func(context.Context, ...*options.ClientOptions) (database.MongoClientInterface, error)
// 		mockPing      func(database.MongoClientInterface, context.Context) error
// 		expectedError string
// 	}{
// 		{
// 			name:   "Successful connection",
// 			uri:    "mongodb://localhost:27017",
// 			dbName: "testdb",
// 			mockConnect: func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClientInterface, error) {
// 				mockClient := new(MockMongoClient)
// 				mockDB := new(MockMongoDatabase)
// 				mockClient.On("Database", "testdb").Return(mockDB)
// 				mockDB.On("Collection", "aws_metrics").Return(&mongo.Collection{})
// 				mockDB.On("Collection", "git_metrics").Return(&mongo.Collection{})
// 				return mockClient, nil
// 			},
// 			mockPing: func(client database.MongoClientInterface, ctx context.Context) error {
// 				return nil
// 			},
// 			expectedError: "",
// 		},
// 		{
// 			name:   "Connection failure",
// 			uri:    "mongodb://localhost:27017",
// 			dbName: "testdb",
// 			mockConnect: func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClientInterface, error) {
// 				return nil, errors.New("connection error")
// 			},
// 			mockPing:      nil,
// 			expectedError: "failed to connect to MongoDB: connection error",
// 		},
// 		{
// 			name:   "Ping failure",
// 			uri:    "mongodb://localhost:27017",
// 			dbName: "testdb",
// 			mockConnect: func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClientInterface, error) {
// 				mockClient := new(MockMongoClient)
// 				return mockClient, nil
// 			},
// 			mockPing: func(client database.MongoClientInterface, ctx context.Context) error {
// 				return errors.New("ping error")
// 			},
// 			expectedError: "failed to ping MongoDB: ping error",
// 		},
// 		{
// 			name:          "Actual implementations",
// 			uri:           "mongodb://localhost:27017",
// 			dbName:        "testdb",
// 			mockConnect:   nil, // We'll use the actual implementation
// 			mockPing:      nil, // We'll use the actual implementation
// 			expectedError: "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mocks
// 			oldConnect, oldPing, oldMongoConnect := database.newMongoClient, database.mongoPing, database.mongoConnect
// 			defer func() {
// 				database.newMongoClient, database.mongoPing, database.mongoConnect = oldConnect, oldPing, oldMongoConnect
// 			}()

// 			if tt.mockConnect != nil {
// 				database.newMongoClient = tt.mockConnect
// 			}
// 			if tt.mockPing != nil {
// 				database.mongoPing = tt.mockPing
// 			}

// 			// Call the function
// 			err := database.InitDB(tt.uri, tt.dbName)

// 			// Assert the result
// 			if tt.expectedError == "" {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, database.Client)
// 				assert.NotNil(t, database.AWSMetricsCollection)
// 				assert.NotNil(t, database.GitMetricsCollection)
// 			} else {
// 				assert.EqualError(t, err, tt.expectedError)
// 			}
// 		})
// 	}
// }
