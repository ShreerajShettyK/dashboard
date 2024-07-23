package helpers

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailTypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

// getLastActivity retrieves the last activity time for an EC2 instance.
func GetLastActivity(ctx context.Context, cloudTrailSvc *cloudtrail.Client, instanceId string) (time.Time, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, -1, 0) // Look back 1 month

	input := &cloudtrail.LookupEventsInput{
		LookupAttributes: []cloudtrailTypes.LookupAttribute{
			{
				AttributeKey:   cloudtrailTypes.LookupAttributeKeyResourceName,
				AttributeValue: aws.String(instanceId),
			},
		},
		StartTime: aws.Time(startTime),
		EndTime:   aws.Time(endTime),
	}

	resp, err := cloudTrailSvc.LookupEvents(ctx, input)
	if err != nil {
		return time.Time{}, err
	}

	if len(resp.Events) > 0 {
		// Events are returned in descending chronological order
		return *resp.Events[0].EventTime, nil
	}

	// If no events found, return the launch time
	return time.Time{}, nil
}
