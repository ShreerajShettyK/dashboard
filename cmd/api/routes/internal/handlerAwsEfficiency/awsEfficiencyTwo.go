package handlerAwsEfficiency

import (
	"context"
	"dashboard/cmd/api/routes/internal/clients"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gorilla/mux"
)

var (
	EC2Client          *ec2.Client
	CostExplorerClient *costexplorer.Client
	CloudTrailClient   *cloudtrail.Client
	CloudWatchClient   *cloudwatch.Client
)

var (
	ListEC2InstancesFunc        = helpers.ListEC2Instances
	FetchEC2InstanceDetailsFunc = helpers.FetchEC2InstanceDetails
	FetchLastActivityFunc       = helpers.FetchLastActivity
	FetchInstanceCostFunc       = helpers.FetchInstanceCost
)

func init() {
	cfg, err := clients.LoadAWSConfig()
	if err != nil {
		return
	}

	EC2Client = clients.NewEC2Client(cfg)
	CostExplorerClient = clients.NewCostExplorerClient(cfg)
	CloudTrailClient = clients.NewCloudTrailClient(cfg)
	CloudWatchClient = clients.NewCloudWatchClient(cfg)
}

func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	services := []string{"ec2", "elb", "rds"}

	data := models.Services{
		Services: services,
	}

	if err := helpers.RenderTemplateFunc(w, data, "aws_billing_Second.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ListServiceInstancesHandler(w http.ResponseWriter, r *http.Request) {
	service := mux.Vars(r)["service"]
	if service == "" {
		http.Error(w, "Service is required", http.StatusBadRequest)
		return
	}

	var instances []string
	var err error

	switch service {
	case "ec2":
		instances, err = ListEC2InstancesFunc(EC2Client)
	case "elb":
		// Implement ELB instance listing
	case "rds":
		// Implement RDS instance listing
	default:
		http.Error(w, "Unsupported service", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	detailedInstances := make([]models.EC2Instance, len(instances))
	for i, instanceId := range instances {
		instanceDetails, err := FetchInstanceDetails(EC2Client, CostExplorerClient, CloudTrailClient, CloudWatchClient, instanceId)
		if err != nil {
			log.Printf("Error fetching instance details for instance %s: %v\n", instanceId, err)
			continue
		}
		detailedInstances[i] = *instanceDetails
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(detailedInstances); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func FetchInstanceDetails(ec2Svc *ec2.Client, ceSvc *costexplorer.Client, cloudTrailSvc *cloudtrail.Client, cloudWatchClient *cloudwatch.Client, instanceId string) (*models.EC2Instance, error) {
	log.Printf("Fetching details for instance ID: %s\n", instanceId)
	instance, err := FetchEC2InstanceDetailsFunc(ec2Svc, instanceId)
	if err != nil {
		return nil, err
	}

	lastActivity, err := FetchLastActivityFunc(context.Background(), cloudWatchClient, instanceId)
	if err != nil {
		log.Printf("Error getting last activity: %v\n", err)
	}

	region := extractRegion(instance)
	daysSinceActivity := int(time.Since(lastActivity).Hours() / 24)

	cost := -1.00
	// cost, err := FetchInstanceCostFunc(ceSvc, instance.InstanceType, region)
	// if err != nil {
	// 	log.Printf("Error getting cost and usage: %v\n", err)
	// 	cost = -1 // Indicate no data
	// }

	location, err := time.LoadLocation("Asia/Kolkata") // Replace with the desired time zone
	if err != nil {
		log.Printf("Error loading location: %v\n", err)
		return nil, err
	}
	lastActivity = lastActivity.In(location)
	log.Printf("cloudwatch last activity: %v", lastActivity.Format("Jan 2, 2006 at 3:04pm"))

	launchedTime := *instance.LaunchTime
	launchedTime = launchedTime.In(location)
	log.Printf("Most recent ec2 launch time: %v", launchedTime.Format("Jan 2, 2006 at 3:04pm"))

	return &models.EC2Instance{
		InstanceId:       instanceId,
		LastActivity:     lastActivity.Format("Jan 2, 2006 at 3:04pm"),
		LastActivityDays: daysSinceActivity,
		Cost:             cost,
		Region:           region,
		InstanceType:     string(instance.InstanceType),
	}, nil
}

func extractRegion(instance *ec2Types.Instance) string {
	if instance.Placement != nil && instance.Placement.AvailabilityZone != nil {
		// Extract region from availability zone
		az := aws.ToString(instance.Placement.AvailabilityZone)
		return az[:len(az)-1] // Remove the last character to get the region
	}
	return ""
}
