package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://task3-shreeraj:YIXZaFDnEmHXC3PS@cluster0.0elhpdy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	gitMetricsCollection := client.Database("task3dashboard").Collection("git_metrics")

	for i := 0; i < 10; i++ {
		doc := bson.M{
			"reponame":       "repo_" + randomString(5),
			"commit_id":      randomString(10),
			"commit_date":    time.Now().AddDate(0, 0, -rand.Intn(100)),
			"commited_by":    "user_" + randomString(5),
			"commit_message": "Fix bug " + randomString(15),
			"files_added":    rand.Intn(10),
			"files_deleted":  rand.Intn(10),
			"files_updated":  rand.Intn(10),
			"lines_added":    rand.Intn(100),
			"lines_updated":  rand.Intn(100),
			"lines_deleted":  rand.Intn(100),
		}

		_, err := gitMetricsCollection.InsertOne(context.Background(), doc)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Git metrics data inserted")
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// below is the code for aws metrics dummy data creation

// package main

// import (
// 	"context"
// 	"log"
// 	"math/rand"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func main() {
// 	clientOptions := options.Client().ApplyURI("mongodb+srv://task3-shreeraj:YIXZaFDnEmHXC3PS@cluster0.0elhpdy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(context.Background())

// 	awsMetricsCollection := client.Database("task3dashboard").Collection("aws_metrics")

// 	services := []string{"ec2", "rds", "elb"}
// 	metrics := []string{"cpu_usage", "disk_space", "memory", "io_reads", "cost"}

// 	for i := 0; i < 10; i++ {
// 		doc := bson.M{
// 			"service_name": services[rand.Intn(len(services))],
// 			"metric_name":  metrics[rand.Intn(len(metrics))],
// 			"value":        rand.Float64() * 100,
// 			"timestamp":    time.Now().AddDate(0, 0, -rand.Intn(100)),
// 		}

// 		_, err := awsMetricsCollection.InsertOne(context.Background(), doc)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	log.Println("AWS metrics data inserted")
// }
