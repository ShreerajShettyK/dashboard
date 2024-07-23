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

func FetchInstanceDetails(ec2Svc *ec2.Client, ceSvc *costexplorer.Client, cloudTrailSvc *cloudtrail.Client, cloudWatchClient *cloudwatch.Client, instanceId string) (*models.EC2Instance, error) {
	log.Printf("Fetching details for instance ID: %s\n", instanceId)
	instance, err := helpers.FetchEC2InstanceDetails(ec2Svc, instanceId)
	if err != nil {
		return nil, err
	}

	// lastActivity, err := helpers.GetLastActivity(context.Background(), cloudTrailSvc, instanceId)
	// if err != nil {
	// 	log.Printf("Error getting last activity: %v\n", err)
	// 	// Fall back to launch time if there's an error
	// 	lastActivity = *instance.LaunchTime
	// }

	lastActivity, err := helpers.FetchLastActivity(context.Background(), cloudWatchClient, instanceId)
	if err != nil {
		log.Printf("Error getting last activity: %v\n", err)
		// Fall back to launch time if there's an error
	}
	// launchedTime := *instance.LaunchTime
	// log.Printf("ec2 launch time is %v", launchedTime.Format("Jan 2, 2006 at 3:04pm"))
	// log.Printf("fetched %v", lastActivity.Format("Jan 2, 2006 at 3:04pm"))

	region := extractRegion(instance)
	daysSinceActivity := int(time.Since(lastActivity).Hours() / 24)

	// cost := -1.00
	cost, err := helpers.FetchInstanceCost(ceSvc, instance.InstanceType, region)
	if err != nil {
		log.Printf("Error getting cost and usage: %v\n", err)
		cost = -1 // Indicate no data
	}

	// Convert the lastActivity time to the desired time zone
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

func ListInstancesHandler(w http.ResponseWriter, r *http.Request) {
	instances, err := helpers.ListEC2Instances(EC2Client)
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

	instanceDetails, err := FetchInstanceDetails(EC2Client, CostExplorerClient, CloudTrailClient, CloudWatchClient, instanceId)
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
	log.Printf("-----------------------------------------------------------------------------------------------------------")
}
