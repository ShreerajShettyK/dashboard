package models

import "time"

type AWSMetric struct {
	ID          int
	ServiceName string
	MetricName  string
	Value       float64
	Timestamp   time.Time
}
