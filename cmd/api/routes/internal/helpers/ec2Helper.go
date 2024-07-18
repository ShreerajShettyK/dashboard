package helpers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ListEC2Instances lists all EC2 instance IDs.
func ListEC2Instances(ec2Svc *ec2.Client) ([]string, error) {
	input := &ec2.DescribeInstancesInput{}
	result, err := ec2Svc.DescribeInstances(context.Background(), input)
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

func FetchEC2InstanceDetails(ec2Svc *ec2.Client, instanceId string) (*ec2Types.Instance, error) {
	instanceInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}
	instanceResult, err := ec2Svc.DescribeInstances(context.Background(), instanceInput)
	if err != nil {
		log.Printf("Error describing instances: %v\n", err)
		return nil, err
	}
	if len(instanceResult.Reservations) > 0 && len(instanceResult.Reservations[0].Instances) > 0 {
		return &instanceResult.Reservations[0].Instances[0], nil
	}
	return nil, nil
}
