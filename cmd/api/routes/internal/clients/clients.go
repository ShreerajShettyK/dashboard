package clients

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// LoadAWSConfig loads the default AWS configuration.
func LoadAWSConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("Error loading AWS config: %v\n", err)
		return aws.Config{}, err
	}
	return cfg, nil
}

// NewEC2Client creates a new EC2 client.
func NewEC2Client(cfg aws.Config) *ec2.Client {
	return ec2.NewFromConfig(cfg)
}

// NewCostExplorerClient creates a new Cost Explorer client.
func NewCostExplorerClient(cfg aws.Config) *costexplorer.Client {
	return costexplorer.NewFromConfig(cfg)
}

// NewCloudTrailClient creates a new CloudTrail client.
func NewCloudTrailClient(cfg aws.Config) *cloudtrail.Client {
	return cloudtrail.NewFromConfig(cfg)
}
