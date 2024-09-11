package service

import (
	"context"
	"errors"
	"os"
	"testing"

	constants "github.com/odhoman/home-devices/internal/constants"
	request "github.com/odhoman/home-devices/internal/request"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDynamoDbApi struct {
	mock.Mock
}

func (m *mockDynamoDbApi) Query(ctx context.Context, input *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *mockDynamoDbApi) PutItem(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *mockDynamoDbApi) GetItem(ctx context.Context, input *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockDynamoDbApi) UpdateItem(ctx context.Context, input *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}
func (m *mockDynamoDbApi) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func TestCreateHomeDevice_Success(t *testing.T) {
	mockEnvVars()
	mockDynamo := new(mockDynamoDbApi)

	mockDynamo.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{}, // Empty result implying no duplicate
	}, nil)

	mockDynamo.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)

	service := HomeDeviceServiceImpl{
		DynamoDbApi: mockDynamo,
	}

	device := request.CreateDeviceRequest{
		MAC:    "00:1B:44:11:3A:B7",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home123",
	}

	result, err := service.CreateHomeDevice(context.Background(), device)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "00:1B:44:11:3A:B7", result.MAC)

	mockDynamo.AssertExpectations(t)
}

func TestCreateHomeDevice_DeviceAlreadyExists(t *testing.T) {
	mockEnvVars()
	mockDynamo := new(mockDynamoDbApi)

	mockDynamo.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{
			{
				"id": &types.AttributeValueMemberS{Value: "existing-device-id"},
			},
		},
	}, nil)

	service := HomeDeviceServiceImpl{
		DynamoDbApi: mockDynamo,
	}

	device := request.CreateDeviceRequest{
		MAC:    "00:1B:44:11:3A:B7",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home123",
	}

	result, err := service.CreateHomeDevice(context.Background(), device)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceAlreadyExistsCode, err.ErrorCode)

	mockDynamo.AssertExpectations(t)
}

func TestCreateHomeDevice_QueryError(t *testing.T) {
	mockEnvVars()
	mockDynamo := new(mockDynamoDbApi)
	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}

	mockDynamo.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{}, errors.New("Query failed"))

	device := request.CreateDeviceRequest{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "Light",
		HomeID: "home1",
	}

	resp, err := service.CreateHomeDevice(context.TODO(), device)

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrGettingDeviceCode, err.ErrorCode)
}

func TestCreateHomeDevice_PutItemError(t *testing.T) {
	mockEnvVars()
	mockDynamo := new(mockDynamoDbApi)
	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}

	mockDynamo.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{}, // Empty result implying no duplicate
	}, nil)

	mockDynamo.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, errors.New("PutItem failed"))

	device := request.CreateDeviceRequest{
		MAC:    "00:11:22:33:44:55",
		Name:   "Test Device",
		Type:   "Light",
		HomeID: "home1",
	}

	resp, err := service.CreateHomeDevice(context.TODO(), device)

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceNotCreatedErrorCode, err.ErrorCode)
}

func TestGetHomeDevice_Success(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberS{Value: "device-id"},
			"mac":        &types.AttributeValueMemberS{Value: "AA:BB:CC:DD:EE:FF"},
			"name":       &types.AttributeValueMemberS{Value: "Test Device"},
			"type":       &types.AttributeValueMemberS{Value: "light"},
			"homeId":     &types.AttributeValueMemberS{Value: "home123"},
			"createdAt":  &types.AttributeValueMemberN{Value: "1626288721"},
			"modifiedAt": &types.AttributeValueMemberN{Value: "1626288721"},
		},
	}, nil)

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	device, err := service.GetHomeDevice(ctx, "device-id")
	assert.NotNil(t, device)
	assert.Nil(t, err)
	assert.Equal(t, "device-id", device.ID)
	assert.Equal(t, "AA:BB:CC:DD:EE:FF", device.MAC)
	assert.Equal(t, "Test Device", device.Name)
	assert.Equal(t, "light", device.Type)
	assert.Equal(t, "home123", device.HomeID)
	assert.Equal(t, int64(1626288721), device.CreatedAt)
	assert.Equal(t, int64(1626288721), device.ModifiedAt)
}

func TestGetHomeDevice_ErrorGettingItem(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{}, errors.New("error getting item"))

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	device, err := service.GetHomeDevice(ctx, "device-id")
	assert.Nil(t, device)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrGettingDeviceCode, err.ErrorCode)
	assert.Equal(t, constants.ErrDeviceNotCreatedErrorMessage, err.ErrorMessage)
}

func TestGetHomeDevice_ItemNotFound(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{Item: nil}, nil)

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	device, err := service.GetHomeDevice(ctx, "device-id")
	assert.Nil(t, device)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceNotFoundCode, err.ErrorCode)
	assert.Equal(t, constants.ErrDeviceNotFoundMessage, err.ErrorMessage)
}

func TestUpdateHomeDevice_Success(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, nil)

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.UpdateHomeDevice(ctx, request.UpdateDeviceRequest{MAC: "AA:BB:CC:DD:EE:FF"}, "device-id")
	assert.Nil(t, err)
}

func TestUpdateHomeDevice_ErrorUpdatingItem(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, errors.New("error updating item"))

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.UpdateHomeDevice(ctx, request.UpdateDeviceRequest{MAC: "AA:BB:CC:DD:EE:FF"}, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrUpdatingDeviceCode, err.ErrorCode)
	assert.Equal(t, constants.ErrUpdatingDeviceMessage, err.ErrorMessage)
}

func TestUpdateHomeDevice_RecordNotFound(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, &types.ConditionalCheckFailedException{})

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.UpdateHomeDevice(ctx, request.UpdateDeviceRequest{MAC: "AA:BB:CC:DD:EE:FF"}, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceNotFoundCode, err.ErrorCode)
	assert.Equal(t, constants.ErrDeviceNotFoundMessage, err.ErrorMessage)
}

func TestUpdateHomeDevice_ErrorBuildingUpdateInput(t *testing.T) {
	clearEnvVars()
	mockDynamo := new(mockDynamoDbApi)
	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.UpdateHomeDevice(ctx, request.UpdateDeviceRequest{MAC: "AA:BB:CC:DD:EE:FF"}, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrGettingConfigCode, err.ErrorCode)
	assert.Equal(t, constants.ErrGettingConfigMessage, err.ErrorMessage)
}

func TestUpdateHomeDevice_NoFieldsToUpdate(t *testing.T) {
	mockEnvVars()
	mockDynamo := new(mockDynamoDbApi)
	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.UpdateHomeDevice(ctx, request.UpdateDeviceRequest{}, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrNoFieldToUpdateCode, err.ErrorCode)
	assert.Equal(t, constants.ErrNoFieldToUpdateMessage, err.ErrorMessage)
}

func TestDeleteHomeDevice_Success(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("DeleteItem", mock.Anything, mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.DeleteHomeDevice(ctx, "device-id")
	assert.Nil(t, err)
}

func TestDeleteHomeDevice_RecordNotFound(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("DeleteItem", mock.Anything, mock.Anything).Return(&dynamodb.DeleteItemOutput{}, &types.ConditionalCheckFailedException{})

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.DeleteHomeDevice(ctx, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceNotFoundCode, err.ErrorCode)
	assert.Equal(t, constants.ErrDeviceNotFoundMessage, err.ErrorMessage)
}

func TestDeleteHomeDevice_ErrorDeletingItem(t *testing.T) {
	mockEnvVars()

	mockDynamo := new(mockDynamoDbApi)
	mockDynamo.On("DeleteItem", mock.Anything, mock.Anything).Return(&dynamodb.DeleteItemOutput{}, errors.New("error deleting item"))

	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.DeleteHomeDevice(ctx, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeletingDeviceCode, err.ErrorCode)
	assert.Equal(t, constants.ErrDeletingDeviceMessage, err.ErrorMessage)
}

func TestDeleteHomeDevice_ErrorGettingTableName(t *testing.T) {
	clearEnvVars()
	mockDynamo := new(mockDynamoDbApi)
	service := HomeDeviceServiceImpl{DynamoDbApi: mockDynamo}
	ctx := context.TODO()

	err := service.DeleteHomeDevice(ctx, "device-id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrGettingConfigCode, err.ErrorCode)
	assert.Equal(t, constants.ErrGettingConfigMessage, err.ErrorMessage)
}

func mockEnvVars() {
	os.Setenv(constants.TableNameHomeDevicesProperty, "table")
	os.Setenv(constants.MacHomeIdIndexNameProperty, "index")
}

func clearEnvVars() {
	os.Setenv(constants.TableNameHomeDevicesProperty, "")
	os.Setenv(constants.MacHomeIdIndexNameProperty, "")
}
