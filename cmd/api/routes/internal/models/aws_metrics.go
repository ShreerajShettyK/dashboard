package models

import "time"

type AWSMetric struct {
	ID          int       `bson:"_id,omitempty"`
	ServiceName string    `bson:"service_name"`
	MetricName  string    `bson:"metric_name"`
	Value       float64   `bson:"value"`
	Timestamp   time.Time `bson:"timestamp"`
}
