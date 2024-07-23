package helpers

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func FetchLastActivity(ctx context.Context, cwSvc *cloudwatch.Client, instanceId string) (time.Time, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, -3, 0) // Look back 3 months

	metrics := []string{"NetworkIn", "NetworkOut", "DiskReadOps", "DiskWriteOps"}
	var latestActivity time.Time

	for _, metric := range metrics {
		timestamp, err := getLastMetricActivity(ctx, cwSvc, instanceId, metric, startTime, endTime)
		if err != nil {
			log.Printf("Error getting %s metric: %v", metric, err)
			continue
		}
		if timestamp.After(latestActivity) {
			latestActivity = timestamp
			log.Printf("New latest activity from %s: %v", metric, latestActivity)
		}
	}

	if latestActivity.IsZero() {
		return time.Time{}, nil
	}

	return latestActivity, nil
}

func getLastMetricActivity(ctx context.Context, cwSvc *cloudwatch.Client, instanceId, metricName string, startTime, endTime time.Time) (time.Time, error) {
	input := &cloudwatch.GetMetricDataInput{
		StartTime: aws.Time(startTime),
		EndTime:   aws.Time(endTime),
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String(metricName),
						Dimensions: []types.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceId),
							},
						},
					},
					Period: aws.Int32(60), // 5 minutes- put 300 (period value is in seconds)
					Stat:   aws.String("Sum"),
				},
				ReturnData: aws.Bool(true),
			},
		},
	}

	result, err := cwSvc.GetMetricData(ctx, input)
	if err != nil {
		log.Printf("Error getting metric data for %s: %v\n", metricName, err)
		return time.Time{}, err
	}

	var latestTimestamp time.Time
	for _, dataResult := range result.MetricDataResults {
		log.Printf("Metric %s: Got %d timestamps and %d values", metricName, len(dataResult.Timestamps), len(dataResult.Values))
		for i := 0; i < len(dataResult.Timestamps); i++ {
			log.Printf("%s - Timestamp: %s, Value: %f\n", metricName, dataResult.Timestamps[i], dataResult.Values[i])
			if dataResult.Values[i] > 0 && dataResult.Timestamps[i].After(latestTimestamp) {
				latestTimestamp = dataResult.Timestamps[i]
			}
		}
	}

	if latestTimestamp.IsZero() {
		log.Printf("No non-zero values found for metric %s", metricName)
		return time.Time{}, nil
	}

	log.Printf("Latest activity for %s: %v", metricName, latestTimestamp)
	return latestTimestamp, nil
}
