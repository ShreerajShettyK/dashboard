package helpers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCloudWatchAPI struct {
	mock.Mock
}

func (m *MockCloudWatchAPI) GetMetricData(ctx context.Context, params *cloudwatch.GetMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.GetMetricDataOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*cloudwatch.GetMetricDataOutput), args.Error(1)
}

func TestFetchLastActivity_Success(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"

	now := time.Now()
	latestTime := now.Add(-1 * time.Hour)

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).Return(&cloudwatch.GetMetricDataOutput{
		MetricDataResults: []types.MetricDataResult{
			{
				Timestamps: []time.Time{latestTime},
				Values:     []float64{1.0},
			},
		},
	}, nil)

	result, err := FetchLastActivity(ctx, mockAPI, instanceId)

	assert.NoError(t, err)
	assert.Equal(t, latestTime, result)
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 4) // One call for each metric
}

func TestFetchLastActivity_NoActivity(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).Return(&cloudwatch.GetMetricDataOutput{
		MetricDataResults: []types.MetricDataResult{
			{
				Timestamps: []time.Time{},
				Values:     []float64{},
			},
		},
	}, nil)

	result, err := FetchLastActivity(ctx, mockAPI, instanceId)

	assert.NoError(t, err)
	assert.True(t, result.IsZero())
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 4) // One call for each metric
}

func TestFetchLastActivity_PartialError(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"

	now := time.Now()
	latestTime := now.Add(-1 * time.Hour)

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).
		Return(&cloudwatch.GetMetricDataOutput{}, errors.New("API error")).Once()

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).
		Return(&cloudwatch.GetMetricDataOutput{
			MetricDataResults: []types.MetricDataResult{
				{
					Timestamps: []time.Time{latestTime},
					Values:     []float64{1.0},
				},
			},
		}, nil).Times(3)

	result, err := FetchLastActivity(ctx, mockAPI, instanceId)

	assert.NoError(t, err)
	assert.Equal(t, latestTime, result)
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 4) // One call for each metric
}

func TestFetchLastActivity_AllErrors(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).
		Return(&cloudwatch.GetMetricDataOutput{}, errors.New("API error"))

	result, err := FetchLastActivity(ctx, mockAPI, instanceId)

	assert.NoError(t, err)
	assert.True(t, result.IsZero())
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 4) // One call for each metric
}

func TestGetLastMetricActivity_Success(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"
	metricName := "NetworkIn"
	startTime := time.Now().Add(-3 * time.Hour)
	endTime := time.Now()

	latestTime := endTime.Add(-1 * time.Hour)

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).Return(&cloudwatch.GetMetricDataOutput{
		MetricDataResults: []types.MetricDataResult{
			{
				Timestamps: []time.Time{latestTime, endTime},
				Values:     []float64{1.0, 0.0},
			},
		},
	}, nil)

	result, err := getLastMetricActivity(ctx, mockAPI, instanceId, metricName, startTime, endTime)

	assert.NoError(t, err)
	assert.Equal(t, latestTime, result)
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 1)
}

func TestGetLastMetricActivity_Error(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"
	metricName := "NetworkIn"
	startTime := time.Now().Add(-3 * time.Hour)
	endTime := time.Now()

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).
		Return(&cloudwatch.GetMetricDataOutput{}, errors.New("API error"))

	result, err := getLastMetricActivity(ctx, mockAPI, instanceId, metricName, startTime, endTime)

	assert.Error(t, err)
	assert.True(t, result.IsZero())
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 1)
}

func TestGetLastMetricActivity_NoActivity(t *testing.T) {
	mockAPI := new(MockCloudWatchAPI)
	ctx := context.Background()
	instanceId := "i-1234567890abcdef0"
	metricName := "NetworkIn"
	startTime := time.Now().Add(-3 * time.Hour)
	endTime := time.Now()

	mockAPI.On("GetMetricData", mock.Anything, mock.Anything, mock.Anything).Return(&cloudwatch.GetMetricDataOutput{
		MetricDataResults: []types.MetricDataResult{
			{
				Timestamps: []time.Time{},
				Values:     []float64{},
			},
		},
	}, nil)

	result, err := getLastMetricActivity(ctx, mockAPI, instanceId, metricName, startTime, endTime)

	assert.NoError(t, err)
	assert.True(t, result.IsZero())
	mockAPI.AssertNumberOfCalls(t, "GetMetricData", 1)
}
