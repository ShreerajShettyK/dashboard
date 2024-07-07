package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var AWSMetricsCollection *mongo.Collection
var GitMetricsCollection *mongo.Collection

// InitDB initializes the MongoDB client and collections.
func InitDB(uri, dbName string) {
	// Set a timeout for the connection context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure client options and connect to MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Assign the global client and collections
	Client = client
	AWSMetricsCollection = client.Database(dbName).Collection("aws_metrics")
	GitMetricsCollection = client.Database(dbName).Collection("git_metrics")

	log.Println("Connected to MongoDB successfully!")
}
