package handlerAwsEfficiency

import (
	"context"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailTypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gorilla/mux"
)

func ListEC2Instances() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.Background(), input)
	if err != nil {
		return nil, err
	}

	var instanceIds []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceIds = append(instanceIds, aws.ToString(instance.InstanceId))
		}
	}
	return instanceIds, nil
}

func getLastActivity(ctx context.Context, cloudTrailSvc *cloudtrail.Client, instanceId string) (time.Time, error) {
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

func FetchInstanceDetails(instanceId string) (*models.EC2Instance, error) {
	log.Printf("Fetching details for instance ID: %s\n", instanceId)
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("Error loading AWS config: %v\n", err)
		return nil, err
	}
	ec2Svc := ec2.NewFromConfig(cfg)
	ceSvc := costexplorer.NewFromConfig(cfg)
	cloudTrailSvc := cloudtrail.NewFromConfig(cfg)

	instanceInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}
	instanceResult, err := ec2Svc.DescribeInstances(context.Background(), instanceInput)
	if err != nil {
		log.Printf("Error describing instances: %v\n", err)
		return nil, err
	}

	var lastActivity time.Time
	var instanceType string
	var region string
	if len(instanceResult.Reservations) > 0 && len(instanceResult.Reservations[0].Instances) > 0 {
		instance := instanceResult.Reservations[0].Instances[0]

		lastActivity, err = getLastActivity(context.Background(), cloudTrailSvc, instanceId)
		if err != nil {
			log.Printf("Error getting last activity: %v\n", err)
			// Fall back to launch time if there's an error
			// lastActivity = *instance.LaunchTime
		}
		// lastActivity = *instance.LaunchTime
		instanceType = string(instance.InstanceType)
		if instance.Placement != nil && instance.Placement.AvailabilityZone != nil {
			// Extract region from availability zone
			az := aws.ToString(instance.Placement.AvailabilityZone)
			region = az[:len(az)-1] // Remove the last character to get the region
		}
	}

	log.Println("Instance type:", instanceType)
	log.Println("Region:", region)

	daysSinceActivity := int(time.Since(lastActivity).Hours() / 24)
	log.Println("Days since activity:", daysSinceActivity)

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
						Values: []string{instanceType},
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

	// Convert the lastActivity time to the desired time zone
	location, err := time.LoadLocation("Asia/Kolkata") // Replace with the desired time zone
	if err != nil {
		log.Printf("Error loading location: %v\n", err)
		return nil, err
	}
	lastActivity = lastActivity.In(location)
	log.Println("Last activity:", lastActivity.Format("Jan 2, 2006 at 3:04pm"))

	log.Println("Sending cost explorer request")
	costResult, err := ceSvc.GetCostAndUsage(context.Background(), costInput)
	if err != nil {
		log.Printf("Error getting cost and usage: %v\n", err)
		// Return partial information without cost
		return &models.EC2Instance{
			InstanceId:       instanceId,
			LastActivity:     lastActivity.Format("Jan 2, 2006 at 3:04pm"),
			LastActivityDays: daysSinceActivity,
			Cost:             -1, // Indicate no data
			Region:           region,
			InstanceType:     instanceType,
		}, nil
	}
	log.Println("Cost result received")

	var totalCost float64
	if len(costResult.ResultsByTime) > 0 {
		for _, resultByTime := range costResult.ResultsByTime {
			for _, group := range resultByTime.Groups {
				cost, _ := strconv.ParseFloat(aws.ToString(group.Metrics["BlendedCost"].Amount), 64)
				totalCost += cost
			}
		}
	}
	log.Printf("Total cost calculated: %f\n", totalCost)

	return &models.EC2Instance{
		InstanceId:       instanceId,
		LastActivity:     lastActivity.Format("Jan 2, 2006 at 3:04pm"),
		LastActivityDays: daysSinceActivity,
		Cost:             totalCost,
		Region:           region,
		InstanceType:     instanceType,
	}, nil
}

func ListInstancesHandler(w http.ResponseWriter, r *http.Request) {
	instances, err := ListEC2Instances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.AWSEfficiencyViewData{
		Instances: instances,
	}

	if err := helpers.RenderTemplateFunc(w, data, "aws_billing.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InstanceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	instanceId := mux.Vars(r)["instance_id"]
	if instanceId == "" {
		http.Error(w, "Instance ID is required", http.StatusBadRequest)
		return
	}

	instanceDetails, err := FetchInstanceDetails(instanceId)
	if err != nil {
		log.Printf("Error fetching instance details: %v\n", err)
		// Return partial information if available
		if instanceDetails != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(instanceDetails)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(instanceDetails); err != nil {
		log.Printf("Error encoding instance details: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Instance details sent for instance ID: %s\n", instanceId)
}
