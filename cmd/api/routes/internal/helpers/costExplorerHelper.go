package helpers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// CostExplorerAPI interface to make the function more testable
type CostExplorerAPI interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

func FetchInstanceCost(ceSvc CostExplorerAPI, instanceType ec2Types.InstanceType, region string) (float64, error) {
	start := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	log.Printf("Cost range: %s to %s\n", start, end)

	costInput := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
		Granularity: types.GranularityDaily,
		Metrics:     []string{"BlendedCost"},
		Filter: &types.Expression{
			And: []types.Expression{
				{
					Dimensions: &types.DimensionValues{
						Key:    types.DimensionService,
						Values: []string{"Amazon Elastic Compute Cloud - Compute"},
					},
				},
				{
					Dimensions: &types.DimensionValues{
						Key:    types.DimensionInstanceType,
						Values: []string{string(instanceType)},
					},
				},
			},
		},
	}

	if region != "" {
		costInput.Filter.And = append(costInput.Filter.And, types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionRegion,
				Values: []string{region},
			},
		})
	}

	log.Println("Sending cost explorer request")
	costResult, err := ceSvc.GetCostAndUsage(context.Background(), costInput)
	if err != nil {
		log.Printf("Error fetching cost data: %v\n", err)
		return 0, fmt.Errorf("failed to fetch cost data: %w", err)
	}
	log.Println("Cost result received")

	var totalCost float64
	if len(costResult.ResultsByTime) > 0 {
		for _, resultByTime := range costResult.ResultsByTime {
			for _, group := range resultByTime.Groups {
				cost, err := strconv.ParseFloat(aws.ToString(group.Metrics["BlendedCost"].Amount), 64)
				if err != nil {
					log.Printf("Error parsing cost amount: %v\n", err)
					return 0, fmt.Errorf("failed to parse cost data: %w", err)
				}
				totalCost += cost
			}
		}
	}
	log.Printf("Total cost calculated: %f\n", totalCost)

	return totalCost, nil
}
