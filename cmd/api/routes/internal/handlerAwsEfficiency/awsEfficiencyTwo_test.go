package handlerAwsEfficiency

import (
	"context"
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock structs
type MockEC2Client struct {
	mock.Mock
}

type MockCostExplorerClient struct {
	mock.Mock
}

type MockCloudTrailClient struct {
	mock.Mock
}

type MockCloudWatchClient struct {
	mock.Mock
}

// Mock helper functions
func mockListEC2Instances(client *ec2.Client) ([]string, error) {
	return []string{"i-1234567890abcdef0"}, nil
}

func mockFetchEC2InstanceDetails(client *ec2.Client, instanceID string) (*ec2Types.Instance, error) {
	return &ec2Types.Instance{
		InstanceId:   aws.String(instanceID),
		InstanceType: ec2Types.InstanceTypeT2Micro,
		Placement: &ec2Types.Placement{
			AvailabilityZone: aws.String("us-west-2a"),
		},
		LaunchTime: aws.Time(time.Now().Add(-24 * time.Hour)),
	}, nil
}

func mockFetchLastActivity(ctx context.Context, client *cloudwatch.Client, instanceID string) (time.Time, error) {
	return time.Now().Add(-12 * time.Hour), nil
}

func mockFetchInstanceCost(client *costexplorer.Client, instanceType string, region string) (float64, error) {
	return 0.5, nil
}

func TestListServicesHandler_Success(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	// Mock the RenderTemplateFunc
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		assert.Equal(t, "aws_billing_Second.html", templateName)
		assert.IsType(t, models.Services{}, data)
		services := data.(models.Services)
		assert.Equal(t, []string{"ec2", "elb", "rds"}, services.Services)
		return nil
	}

	req, err := http.NewRequest("GET", "/services", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListServicesHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListServicesHandler_RenderError(t *testing.T) {
	// Save the original function and defer its restoration
	originalRenderTemplateFunc := helpers.RenderTemplateFunc
	defer func() { helpers.RenderTemplateFunc = originalRenderTemplateFunc }()

	// Mock the RenderTemplateFunc
	helpers.RenderTemplateFunc = func(w http.ResponseWriter, data interface{}, templateName string) error {
		return errors.New("render error")
	}

	req, err := http.NewRequest("GET", "/services", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListServicesHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "render error\n", rr.Body.String())
}

// func TestFetchInstanceDetails(t *testing.T) {
// 	// Setup mock clients
// 	mockEC2 := new(MockEC2Client)
// 	mockCE := new(MockCostExplorerClient)
// 	mockCT := new(MockCloudTrailClient)
// 	mockCW := new(MockCloudWatchClient)

// 	// Save original functions and restore after test
// 	originalFetchEC2InstanceDetailsFunc := FetchEC2InstanceDetailsFunc
// 	originalFetchLastActivityFunc := FetchLastActivityFunc
// 	originalFetchInstanceCostFunc := FetchInstanceCostFunc
// 	defer func() {
// 		FetchEC2InstanceDetailsFunc = originalFetchEC2InstanceDetailsFunc
// 		FetchLastActivityFunc = originalFetchLastActivityFunc
// 		FetchInstanceCostFunc = originalFetchInstanceCostFunc
// 	}()

// 	instanceID := "i-1234567890abcdef0"

// 	// Positive case
// 	t.Run("Successful fetch", func(t *testing.T) {
// 		FetchEC2InstanceDetailsFunc = mockFetchEC2InstanceDetails
// 		FetchLastActivityFunc = mockFetchLastActivity
// 		FetchInstanceCostFunc = mockFetchInstanceCost

// 		instance, err := FetchInstanceDetails(mockEC2, mockCE, mockCT, mockCW, instanceID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, instance)
// 		assert.Equal(t, instanceID, instance.InstanceId)
// 		assert.Equal(t, "us-west-2", instance.Region)
// 		assert.Equal(t, "t2.micro", instance.InstanceType)
// 		assert.InDelta(t, 0.5, instance.Cost, 0.01)
// 	})

// 	// Negative cases
// 	t.Run("Failed EC2 instance fetch", func(t *testing.T) {
// 		FetchEC2InstanceDetailsFunc = func(client *ec2.Client, instanceID string) (*ec2Types.Instance, error) {
// 			return nil, errors.New("EC2 fetch error")
// 		}

// 		instance, err := FetchInstanceDetails(mockEC2, mockCE, mockCT, mockCW, instanceID)

// 		assert.Error(t, err)
// 		assert.Nil(t, instance)
// 		assert.Contains(t, err.Error(), "EC2 fetch error")
// 	})

// 	t.Run("Failed last activity fetch", func(t *testing.T) {
// 		FetchEC2InstanceDetailsFunc = mockFetchEC2InstanceDetails
// 		FetchLastActivityFunc = func(ctx context.Context, client *cloudwatch.Client, instanceID string) (time.Time, error) {
// 			return time.Time{}, errors.New("Last activity fetch error")
// 		}
// 		FetchInstanceCostFunc = mockFetchInstanceCost

// 		instance, err := FetchInstanceDetails(mockEC2, mockCE, mockCT, mockCW, instanceID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, instance)
// 		assert.Equal(t, "Jan 1, 0001 at 12:00am", instance.LastActivity)
// 	})

// 	t.Run("Failed instance cost fetch", func(t *testing.T) {
// 		FetchEC2InstanceDetailsFunc = mockFetchEC2InstanceDetails
// 		FetchLastActivityFunc = mockFetchLastActivity
// 		FetchInstanceCostFunc = func(client *costexplorer.Client, instanceType string, region string) (float64, error) {
// 			return 0, errors.New("Cost fetch error")
// 		}

// 		instance, err := FetchInstanceDetails(mockEC2, mockCE, mockCT, mockCW, instanceID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, instance)
// 		assert.Equal(t, -1.0, instance.Cost)
// 	})
// }

func TestExtractRegion(t *testing.T) {
	testCases := []struct {
		name           string
		instance       *ec2Types.Instance
		expectedRegion string
	}{
		{
			name: "Valid AZ",
			instance: &ec2Types.Instance{
				Placement: &ec2Types.Placement{
					AvailabilityZone: aws.String("us-west-2a"),
				},
			},
			expectedRegion: "us-west-2",
		},
		{
			name:           "Nil Placement",
			instance:       &ec2Types.Instance{},
			expectedRegion: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			region := extractRegion(tc.instance)
			assert.Equal(t, tc.expectedRegion, region)
		})
	}
}
