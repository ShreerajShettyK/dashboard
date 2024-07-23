package handlerAwsEfficiency

import (
	"dashboard/cmd/api/routes/internal/helpers"
	"dashboard/cmd/api/routes/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

// Mock clients
// type MockEC2Client struct {
// 	mock.Mock
// }

// type MockCostExplorerClient struct {
// 	mock.Mock
// }

// type MockCloudTrailClient struct {
// 	mock.Mock
// }

// type MockCloudWatchClient struct {
// 	mock.Mock
// }

// // Mock helper functions
// func mockListEC2Instances(client *ec2.Client) ([]string, error) {
// 	return []string{"i-1234567890abcdef0"}, nil
// }

// func mockFetchEC2InstanceDetails(client *ec2.Client, instanceID string) (*ec2Types.Instance, error) {
// 	return &ec2Types.Instance{
// 		InstanceId:   aws.String(instanceID),
// 		InstanceType: ec2Types.InstanceTypeT2Micro,
// 		Placement: &ec2Types.Placement{
// 			AvailabilityZone: aws.String("us-west-2a"),
// 		},
// 		LaunchTime: aws.Time(time.Now().Add(-24 * time.Hour)),
// 	}, nil
// }

// func mockFetchLastActivity(ctx context.Context, client *cloudwatch.Client, instanceID string) (time.Time, error) {
// 	return time.Now().Add(-12 * time.Hour), nil
// }

// func mockFetchInstanceCost(client *costexplorer.Client, instanceType string, region string) (float64, error) {
// 	return 0.5, nil
// }

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

// func TestListServiceInstancesHandler(t *testing.T) {
// 	// Setup
// 	EC2Client = &ec2.Client{}
// 	CostExplorerClient = &costexplorer.Client{}
// 	CloudTrailClient = &cloudtrail.Client{}
// 	CloudWatchClient = &cloudwatch.Client{}

// 	// Test cases
// 	testCases := []struct {
// 		name           string
// 		service        string
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{"Valid EC2 Service", "ec2", http.StatusOK, `[{"instance_id":"i-1234567890abcdef0"`},
// 		{"Invalid Service", "invalid", http.StatusBadRequest, "Unsupported service"},
// 		{"Empty Service", "", http.StatusBadRequest, "Service is required"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			req, err := http.NewRequest("GET", "/services/"+tc.service+"/instances", nil)
// 			assert.NoError(t, err)

// 			rr := httptest.NewRecorder()
// 			router := mux.NewRouter()
// 			router.HandleFunc("/services/{service}/instances", ListServiceInstancesHandler)

// 			router.ServeHTTP(rr, req)

// 			assert.Equal(t, tc.expectedStatus, rr.Code)
// 			assert.Contains(t, rr.Body.String(), tc.expectedBody)
// 		})
// 	}
// }

// func TestFetchInstanceDetails(t *testing.T) {
// 	// Setup mock clients
// 	mockEC2 := new(MockEC2Client)
// 	mockCE := new(MockCostExplorerClient)
// 	mockCT := new(MockCloudTrailClient)
// 	mockCW := new(MockCloudWatchClient)

// 	instanceID := "i-1234567890abcdef0"

// 	// Positive case
// 	t.Run("Successful fetch", func(t *testing.T) {
// 		mockEC2.On("DescribeInstances", mock.Anything, mock.Anything).Return(&ec2.DescribeInstancesOutput{
// 			Reservations: []ec2Types.Reservation{
// 				{
// 					Instances: []ec2Types.Instance{
// 						{
// 							InstanceId:   aws.String(instanceID),
// 							InstanceType: ec2Types.InstanceTypeT2Micro,
// 							Placement: &ec2Types.Placement{
// 								AvailabilityZone: aws.String("us-west-2a"),
// 							},
// 							LaunchTime: aws.Time(time.Now().Add(-24 * time.Hour)),
// 						},
// 					},
// 				},
// 			},
// 		}, nil)

// 		mockCW.On("GetMetricData", mock.Anything, mock.Anything).Return(&cloudwatch.GetMetricDataOutput{
// 			MetricDataResults: []cloudwatch.MetricDataResult{
// 				{
// 					Timestamps: []time.Time{time.Now().Add(-12 * time.Hour)},
// 					Values:     []float64{1.0},
// 				},
// 			},
// 		}, nil)

// 		mockCE.On("GetCostAndUsage", mock.Anything, mock.Anything).Return(&costexplorer.GetCostAndUsageOutput{
// 			ResultsByTime: []costexplorer.ResultByTime{
// 				{
// 					Total: map[string]costexplorer.MetricValue{
// 						"UnblendedCost": {
// 							Amount: aws.String("0.5"),
// 						},
// 					},
// 				},
// 			},
// 		}, nil)

// 		instance, err := FetchInstanceDetails(mockEC2, mockCE, mockCT, mockCW, instanceID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, instance)
// 		assert.Equal(t, instanceID, instance.InstanceId)
// 		assert.Equal(t, "us-west-2", instance.Region)
// 		assert.Equal(t, "t2.micro", instance.InstanceType)
// 		assert.InDelta(t, 0.5, instance.Cost, 0.01)
// 	})

// 	// Negative case
// 	t.Run("Failed fetch", func(t *testing.T) {
// 		mockEC2.On("DescribeInstances", mock.Anything, mock.Anything).Return(nil, assert.AnError)

// 		instance, err := FetchInstanceDetails(mockEC2, mockCE, mockCT, mockCW, instanceID)

// 		assert.Error(t, err)
// 		assert.Nil(t, instance)
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
