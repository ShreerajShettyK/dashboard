// package models

// import "time"

// type AWSMetric struct {
// 	ServiceName string    `bson:"service_name"`
// 	Date        time.Time `bson:"date"`
// 	CPUUsage    float64   `bson:"cpu_usage"`
// 	DiskSpace   float64   `bson:"disk_space"`
// 	Memory      float64   `bson:"memory"`
// 	IOReads     float64   `bson:"io_reads"`
// 	IOWrites    float64   `bson:"io_writes"`
// 	NetworkIn   float64   `bson:"network_in"`
// 	NetworkOut  float64   `bson:"network_out"`
// 	Cost        float64   `bson:"cost"`
// }

// // GitMetricsViewData holds the data for rendering the Git metrics view
// type AwsMetricsViewData struct {
// 	Metrics  []AWSMetric
// 	Services []string
// }

package models

import "time"

type AWSMetric struct {
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

type AwsMetricsViewData struct {
	Metrics     []AWSMetric
	Services    []string
	ServiceName string
	StartDate   string
	EndDate     string
}
