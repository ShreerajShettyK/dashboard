package helpers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCostExplorerClient struct {
	mock.Mock
}

func (m *MockCostExplorerClient) GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*costexplorer.GetCostAndUsageOutput), args.Error(1)
}

func TestFetchInstanceCost(t *testing.T) {
	mockClient := new(MockCostExplorerClient)

	tests := []struct {
		name         string
		instanceType ec2Types.InstanceType
		region       string
		setupMock    func()
		expected     float64
		expectedErr  error
	}{
		{
			name:         "Successful cost fetch",
			instanceType: ec2Types.InstanceTypeT2Micro,
			region:       "us-west-2",
			setupMock: func() {
				mockClient.On("GetCostAndUsage", mock.Anything, mock.AnythingOfType("*costexplorer.GetCostAndUsageInput"), mock.Anything).Return(
					&costexplorer.GetCostAndUsageOutput{
						ResultsByTime: []types.ResultByTime{
							{
								Groups: []types.Group{
									{
										Metrics: map[string]types.MetricValue{
											"BlendedCost": {Amount: aws.String("10.5")},
										},
									},
								},
							},
						},
					},
					nil,
				)
			},
			expected:    10.5,
			expectedErr: nil,
		},
		{
			name:         "Empty region",
			instanceType: ec2Types.InstanceTypeT2Micro,
			region:       "",
			setupMock: func() {
				mockClient.On("GetCostAndUsage", mock.Anything, mock.AnythingOfType("*costexplorer.GetCostAndUsageInput"), mock.Anything).Return(
					&costexplorer.GetCostAndUsageOutput{
						ResultsByTime: []types.ResultByTime{
							{
								Groups: []types.Group{
									{
										Metrics: map[string]types.MetricValue{
											"BlendedCost": {Amount: aws.String("10.5")},
										},
									},
								},
							},
						},
					},
					nil,
				)
			},
			expected:    10.5,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := FetchInstanceCost(mockClient, tt.instanceType, tt.region)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFetchInstanceCostSecond(t *testing.T) {
	mockClient := new(MockCostExplorerClient)
	tests := []struct {
		name         string
		instanceType ec2Types.InstanceType
		region       string
		setupMock    func()
		expected     float64
		expectedErr  string
	}{
		{
			name:         "GetCostAndUsage error",
			instanceType: ec2Types.InstanceTypeT2Micro,
			region:       "us-west-2",
			setupMock: func() {
				mockClient.On("GetCostAndUsage", mock.Anything, mock.AnythingOfType("*costexplorer.GetCostAndUsageInput"), mock.Anything).Return(
					(*costexplorer.GetCostAndUsageOutput)(nil),
					errors.New("API error"),
				)
			},
			expected:    0,
			expectedErr: "failed to fetch cost data: API error",
		},
		{
			name:         "Invalid cost amount",
			instanceType: ec2Types.InstanceTypeT2Micro,
			region:       "us-west-2",
			setupMock: func() {
				mockClient.On("GetCostAndUsage", mock.Anything, mock.AnythingOfType("*costexplorer.GetCostAndUsageInput"), mock.Anything).Return(
					&costexplorer.GetCostAndUsageOutput{
						ResultsByTime: []types.ResultByTime{
							{
								Groups: []types.Group{
									{
										Metrics: map[string]types.MetricValue{
											"BlendedCost": {Amount: aws.String("invalid")},
										},
									},
								},
							},
						},
					},
					nil,
				)
			},
			expected:    0,
			expectedErr: "failed to parse cost data: strconv.ParseFloat: parsing \"invalid\": invalid syntax",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient = new(MockCostExplorerClient) // Reset mock for each test
			tt.setupMock()
			result, err := FetchInstanceCost(mockClient, tt.instanceType, tt.region)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestFetchInstanceCost_DateRange(t *testing.T) {
	mockClient := new(MockCostExplorerClient)

	mockClient.On("GetCostAndUsage", mock.Anything, mock.AnythingOfType("*costexplorer.GetCostAndUsageInput"), mock.Anything).Return(
		&costexplorer.GetCostAndUsageOutput{
			ResultsByTime: []types.ResultByTime{},
		},
		nil,
	)

	_, err := FetchInstanceCost(mockClient, ec2Types.InstanceTypeT2Micro, "us-west-2")

	assert.NoError(t, err)

	capturedInput := mockClient.Calls[0].Arguments.Get(1).(*costexplorer.GetCostAndUsageInput)

	// Check if the date range is correct (30 days)
	start, _ := time.Parse("2006-01-02", *capturedInput.TimePeriod.Start)
	end, _ := time.Parse("2006-01-02", *capturedInput.TimePeriod.End)
	assert.Equal(t, 30, int(end.Sub(start).Hours()/24))

	mockClient.AssertExpectations(t)
}
