package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Metric struct {
	ServiceName string    `bson:"service_name"`
	Date        time.Time `bson:"date"`
	CPUUsage    float64   `bson:"cpu_usage"`
	DiskSpace   float64   `bson:"disk_space"`
	Memory      float64   `bson:"memory"`
	IOReads     float64   `bson:"io_reads"`
	IOWrites    float64   `bson:"io_writes"`
	NetworkIn   float64   `bson:"network_in"`
	NetworkOut  float64   `bson:"network_out"`
	Cost        float64   `bson:"cost"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://task3-shreeraj:YIXZaFDnEmHXC3PS@cluster0.0elhpdy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	awsMetricsCollection := client.Database("dashboard").Collection("aws_metrics")
	services := []string{"ec2", "rds", "elb"}

	// Generate dummy data for the month of July 2023
	var metrics []interface{}
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) //year month date
	for i := 0; i < 30; i++ {
		date := startDate.AddDate(0, 0, i)
		metrics = append(metrics, Metric{
			ServiceName: services[rand.Intn(len(services))],
			Date:        date,
			CPUUsage:    50 + float64(i), // Dummy values
			DiskSpace:   500,
			Memory:      16000,
			IOReads:     1000 + float64(i*10),
			IOWrites:    500 + float64(i*5),
			NetworkIn:   200 + float64(i*2),
			NetworkOut:  100 + float64(i),
			Cost:        50 + float64(i),
		})
	}

	_, err = awsMetricsCollection.InsertMany(ctx, metrics)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dummy data inserted successfully")
}
