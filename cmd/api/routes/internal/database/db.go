// 86.7%

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client               *mongo.Client
	AWSMetricsCollection *mongo.Collection
	GitMetricsCollection *mongo.Collection
	mongoConnect         = mongo.Connect
	mongoPing            = func(client *mongo.Client, ctx context.Context) error { return client.Ping(ctx, nil) }
	newMongoClient       = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
		return mongoConnect(ctx, opts...)
	}
)

func InitDB(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := newMongoClient(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = mongoPing(client, ctx)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	Client = client
	AWSMetricsCollection = client.Database(dbName).Collection("aws_metrics")
	GitMetricsCollection = client.Database(dbName).Collection("git_metrics")

	log.Println("Connected to MongoDB successfully!")
	return nil
}

// package database

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"go.mongodb.org/mongo-driver/mongo/readpref"
// )

// type MongoClientInterface interface {
// 	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
// 	Ping(ctx context.Context, rp *readpref.ReadPref) error
// }

// var (
// 	Client               MongoClientInterface
// 	AWSMetricsCollection *mongo.Collection
// 	GitMetricsCollection *mongo.Collection

// 	mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (MongoClientInterface, error) {
// 		client, err := mongo.Connect(ctx, opts...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return client, nil
// 	}
// 	mongoPing = func(client MongoClientInterface, ctx context.Context) error {
// 		return client.Ping(ctx, nil)
// 	}
// 	newMongoClient = func(ctx context.Context, opts ...*options.ClientOptions) (MongoClientInterface, error) {
// 		return mongoConnect(ctx, opts...)
// 	}
// )

// func InitDB(uri, dbName string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	clientOptions := options.Client().ApplyURI(uri)
// 	client, err := newMongoClient(ctx, clientOptions)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to MongoDB: %v", err)
// 	}

// 	err = mongoPing(client, ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to ping MongoDB: %v", err)
// 	}

// 	Client = client
// 	AWSMetricsCollection = client.Database(dbName).Collection("aws_metrics")
// 	GitMetricsCollection = client.Database(dbName).Collection("git_metrics")

// 	log.Println("Connected to MongoDB successfully!")
// 	return nil
// }
