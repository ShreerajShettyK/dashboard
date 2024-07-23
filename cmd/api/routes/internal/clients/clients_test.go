package clients

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
)

// Mocking the LoadDefaultConfig function
var (
	originalLoadDefaultConfig = loadDefaultConfig
)

func mockLoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return aws.Config{}, nil
}

func mockLoadDefaultConfigError(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return aws.Config{}, errors.New("load config error")
}

func TestLoadAWSConfig(t *testing.T) {
	loadDefaultConfig = mockLoadDefaultConfig
	defer func() { loadDefaultConfig = originalLoadDefaultConfig }()

	cfg, err := LoadAWSConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestLoadAWSConfig_Error(t *testing.T) {
	loadDefaultConfig = mockLoadDefaultConfigError
	defer func() { loadDefaultConfig = originalLoadDefaultConfig }()

	cfg, err := LoadAWSConfig()
	assert.Error(t, err)
	assert.Equal(t, aws.Config{}, cfg)
}

func TestNewEC2Client(t *testing.T) {
	cfg := aws.Config{}
	client := NewEC2Client(cfg)
	assert.NotNil(t, client)
	assert.IsType(t, &ec2.Client{}, client)
}

func TestNewCostExplorerClient(t *testing.T) {
	cfg := aws.Config{}
	client := NewCostExplorerClient(cfg)
	assert.NotNil(t, client)
	assert.IsType(t, &costexplorer.Client{}, client)
}

func TestNewCloudTrailClient(t *testing.T) {
	cfg := aws.Config{}
	client := NewCloudTrailClient(cfg)
	assert.NotNil(t, client)
	assert.IsType(t, &cloudtrail.Client{}, client)
}

func TestNewCloudWatchClient(t *testing.T) {
	cfg := aws.Config{}
	client := NewCloudWatchClient(cfg)
	assert.NotNil(t, client)
	assert.IsType(t, &cloudwatch.Client{}, client)
}
