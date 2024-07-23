package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2Client is a mock of the EC2 client
type MockEC2Client struct {
	mock.Mock
}

func (m *MockEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
}

func TestListEC2Instances(t *testing.T) {
	mockClient := new(MockEC2Client)

	tests := []struct {
		name        string
		setupMock   func()
		expectedIDs []string
		expectedErr error
	}{
		{
			name: "Successful listing",
			setupMock: func() {
				mockClient.On("DescribeInstances", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.Anything).Return(
					&ec2.DescribeInstancesOutput{
						Reservations: []ec2Types.Reservation{
							{
								Instances: []ec2Types.Instance{
									{InstanceId: aws.String("i-1234567890abcdef0")},
									{InstanceId: aws.String("i-0987654321fedcba0")},
								},
							},
						},
					},
					nil,
				)
			},
			expectedIDs: []string{"i-1234567890abcdef0", "i-0987654321fedcba0"},
			expectedErr: nil,
		},
		{
			name: "API error",
			setupMock: func() {
				mockClient.On("DescribeInstances", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.Anything).Return(
					(*ec2.DescribeInstancesOutput)(nil),
					errors.New("API error"),
				)
			},
			expectedIDs: nil,
			expectedErr: errors.New("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient = new(MockEC2Client)
			tt.setupMock()

			ids, err := ListEC2Instances(mockClient)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIDs, ids)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFetchEC2InstanceDetails(t *testing.T) {
	mockClient := new(MockEC2Client)

	tests := []struct {
		name             string
		instanceID       string
		setupMock        func()
		expectedInstance *ec2Types.Instance
		expectedErr      error
	}{
		{
			name:       "Successful fetch",
			instanceID: "i-1234567890abcdef0",
			setupMock: func() {
				mockClient.On("DescribeInstances", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.Anything).Return(
					&ec2.DescribeInstancesOutput{
						Reservations: []ec2Types.Reservation{
							{
								Instances: []ec2Types.Instance{
									{InstanceId: aws.String("i-1234567890abcdef0")},
								},
							},
						},
					},
					nil,
				)
			},
			expectedInstance: &ec2Types.Instance{InstanceId: aws.String("i-1234567890abcdef0")},
			expectedErr:      nil,
		},
		{
			name:       "API error",
			instanceID: "i-1234567890abcdef0",
			setupMock: func() {
				mockClient.On("DescribeInstances", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.Anything).Return(
					(*ec2.DescribeInstancesOutput)(nil),
					errors.New("API error"),
				)
			},
			expectedInstance: nil,
			expectedErr:      errors.New("API error"),
		},
		{
			name:       "Instance not found",
			instanceID: "i-1234567890abcdef0",
			setupMock: func() {
				mockClient.On("DescribeInstances", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.Anything).Return(
					&ec2.DescribeInstancesOutput{
						Reservations: []ec2Types.Reservation{},
					},
					nil,
				)
			},
			expectedInstance: nil,
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient = new(MockEC2Client)
			tt.setupMock()

			instance, err := FetchEC2InstanceDetails(mockClient, tt.instanceID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInstance, instance)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
