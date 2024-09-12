package main

import (
	"context"
	"testing"

	"github.com/odhoman/home-devices/internal/dao"
	"github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"
	hdResponse "github.com/odhoman/home-devices/internal/response"
	hDService "github.com/odhoman/home-devices/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mock.RunTestMain(m)
}

func TestDeleteHomeDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "34:76:2B:BB:EE:45",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	ctx := context.Background()
	svc := mock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	deviceCreated := CreateHomeDeviceForTesting(t, ctx, homeDeviceServiceImpl, request)

	response, err := HandleRequest(context.Background(), deviceCreated.ID, homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("Unexpected error running TestDeleteHomeDevice_Success. Error: %v", err.Error())
	}

	assert.NotNil(t, response)
	assert.Contains(t, response.Body, "Device deleted")
	assert.Equal(t, response.StatusCode, 200)
}

func TestDeleteHomeDevice_DeviceNotFound(t *testing.T) {

	svc := mock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	response, err := HandleRequest(context.Background(), "fakeId", homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("Unexpected error running TestDeleteHomeDevice. Error: %v", err.Error())
	}

	assert.Contains(t, response.Body, "Device Not Found")
	assert.Equal(t, response.StatusCode, 404)
}

func CreateHomeDeviceForTesting(t *testing.T, ctx context.Context, service hDService.HomeDeviceService, device hDRequest.CreateDeviceRequest) *hdResponse.HomdeDeviceResponse {

	deviceCreated, err := service.CreateHomeDevice(ctx, device)

	if err != nil {
		t.Fatalf("Unexpected error creating new Device for testing. Error: %v", err.ErrorCode)
		return nil
	}

	return deviceCreated
}
