package main

import (
	"context"
	"testing"

	"github.com/odhoman/home-devices/internal/dao"
	"github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"
	hDResponse "github.com/odhoman/home-devices/internal/response"
	hDService "github.com/odhoman/home-devices/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mock.RunTestMain(m)
}

func TestHandleUpdateDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "00:1A:2B:3C:4D:5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	ctx := context.Background()
	svc := mock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})
	deviceCreated := CreateHomeDeviceForTesting(t, ctx, homeDeviceServiceImpl, request)

	id := deviceCreated.ID

	updateRequest := hDRequest.UpdateDeviceRequest{
		Type: "alarm",
	}

	response, err := HandleRequest(ctx, updateRequest, deviceCreated.ID, homeDeviceServiceImpl)

	deviceReturned := GetHomeDeviceForTesting(t, ctx, homeDeviceServiceImpl, id)

	assert.NoError(t, err)
	assert.Equal(t, deviceReturned.ID, id)
	assert.Equal(t, deviceReturned.Type, "alarm")
	assert.NotNil(t, response)
	assert.Contains(t, response.Body, "Device updated")
	assert.Equal(t, response.StatusCode, 201)

}

func TestHandleUpdateDevice_DeviceNotFound(t *testing.T) {

	ctx := context.Background()
	svc := mock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	updateRequest := hDRequest.UpdateDeviceRequest{
		Type: "alarm",
	}

	response, err := HandleRequest(ctx, updateRequest, "fakeId", homeDeviceServiceImpl)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Contains(t, response.Body, "Device Not Found")
	assert.Equal(t, response.StatusCode, 404)
}

func CreateHomeDeviceForTesting(t *testing.T, ctx context.Context, service hDService.HomeDeviceService, device hDRequest.CreateDeviceRequest) *hDResponse.HomdeDeviceResponse {

	deviceCreated, err := service.CreateHomeDevice(ctx, device)

	if err != nil {
		t.Fatalf("Unexpected error creating new Device for testing. Error: %v", err.ErrorCode)
		return nil
	}

	return deviceCreated
}

func GetHomeDeviceForTesting(t *testing.T, ctx context.Context, service hDService.HomeDeviceService, id string) *hDResponse.HomdeDeviceResponse {

	deviceReturned, err := service.GetHomeDevice(ctx, id)

	if err != nil {
		t.Fatalf("Unexpected error creating new Device for testing. Error: %v", err.ErrorCode)
		return nil
	}

	return deviceReturned
}
