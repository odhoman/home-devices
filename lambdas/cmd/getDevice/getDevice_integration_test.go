package main

import (
	"context"
	"testing"

	"github.com/odhoman/home-devices/internal/dao"
	hdMock "github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"
	hdResponse "github.com/odhoman/home-devices/internal/response"
	hDService "github.com/odhoman/home-devices/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	hdMock.RunTestMain(m)
}

func TestGetDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "98:76:AA:BB:EE:45",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	ctx := context.Background()
	svc := hdMock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	deviceCreated := createHomeDeviceForTesting(t, ctx, homeDeviceServiceImpl, request)

	response, err := HandleRequest(context.Background(), deviceCreated.ID, homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("Unexpected error running TestGetDevice_Success. Error: %v", err.Error())
	}

	assert.NotNil(t, response)
	assert.Contains(t, response.Body, "98:76:AA:BB:EE:45")
	assert.Contains(t, response.Body, "Living Room Light")
	assert.Contains(t, response.Body, "home12122")
	assert.Equal(t, response.StatusCode, 200)
}

func TestGetDevice_DeviceNotFound(t *testing.T) {

	svc := hdMock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	response, err := HandleRequest(context.Background(), "fakeId", homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("Unexpected error running TestGetDevice_DeviceNotFound. Error: %v", err.Error())
	}

	assert.Contains(t, response.Body, "Device Not Found")
	assert.Equal(t, response.StatusCode, 404)
}

func createHomeDeviceForTesting(t *testing.T, ctx context.Context, service hDService.HomeDeviceService, device hDRequest.CreateDeviceRequest) *hdResponse.HomdeDeviceResponse {

	deviceCreated, err := service.CreateHomeDevice(ctx, device)

	if err != nil {
		t.Fatalf("Unexpected error creating new Device for testing. Error: %v", err.ErrorCode)
		return nil
	}

	return deviceCreated
}
