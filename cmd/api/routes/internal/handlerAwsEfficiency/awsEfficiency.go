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
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gorilla/mux"
)

type EC2Instance struct {
	InstanceId       string  `json:"InstanceId"`
	LastActivity     string  `json:"LastActivity"`
	LastActivityDays int     `json:"LastActivityDays"`
	Cost             float64 `json:"Cost"`
}

func ListEC2Instances() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var instanceIds []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
	}
	return instanceIds, nil
}

func FetchInstanceDetails(instanceId string) (*EC2Instance, error) {
	log.Printf("Fetching details for instance ID: %s\n", instanceId)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading AWS config: %v\n", err)
		return nil, err
	}
	ec2Svc := ec2.NewFromConfig(cfg)
	ceSvc := costexplorer.NewFromConfig(cfg)

	instanceInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}
	instanceResult, err := ec2Svc.DescribeInstances(context.TODO(), instanceInput)
	if err != nil {
		log.Printf("Error describing instances: %v\n", err)
		return nil, err
	}

	var lastActivity time.Time
	if len(instanceResult.Reservations) > 0 && len(instanceResult.Reservations[0].Instances) > 0 {
		lastActivity = *instanceResult.Reservations[0].Instances[0].LaunchTime
	}
	log.Println("Last activity:", lastActivity)

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
			Dimensions: &types.DimensionValues{
				Key:    types.Dimension("INSTANCE_ID"),
				Values: []string{instanceId},
			},
		},
	}

	log.Println("Sending cost explorer request")
	costResult, err := ceSvc.GetCostAndUsage(context.TODO(), costInput)
	if err != nil {
		log.Printf("Error getting cost and usage: %v\n", err)
		// Return partial information without cost
		return &EC2Instance{
			InstanceId:       instanceId,
			LastActivity:     lastActivity.Format("Jan 2, 2006 at 3:04pm"),
			LastActivityDays: daysSinceActivity,
			Cost:             0, // or you could use -1 to indicate no data
		}, nil
	}

	log.Println("Cost result received")

	var totalCost float64
	if len(costResult.ResultsByTime) > 0 {
		for _, resultByTime := range costResult.ResultsByTime {
			for _, group := range resultByTime.Groups {
				cost, _ := strconv.ParseFloat(*group.Metrics["BlendedCost"].Amount, 64)
				totalCost += cost
			}
		}
	}
	log.Printf("Total cost calculated: %f\n", totalCost)

	return &EC2Instance{
		InstanceId:       instanceId,
		LastActivity:     lastActivity.Format("Jan 2, 2006 at 3:04pm"),
		LastActivityDays: daysSinceActivity,
		Cost:             totalCost,
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
