package main

import (
	"context"
	"testing"

	"github.com/odhoman/home-devices/internal/dao"
	"github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"
	hDService "github.com/odhoman/home-devices/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mock.RunTestMain(m)
}

func TestCreateHomeDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "00:1A:2B:3C:4D:5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	svc := mock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	response, err := HandleRequest(context.Background(), request, homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("Unexpected error running TestCreateHomeDevice_Success. Error: %v", err.Error())
	}

	assert.NotNil(t, response)
	assert.Contains(t, response.Body, "00:1A:2B:3C:4D:5E")
	assert.Equal(t, response.StatusCode, 201)
}

func TestCreateHomeDevice_DeviceAlreadyExist(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "0A:1B:2C:3D:4E:5F",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	svc := mock.GetDynamoConnectionTestFromEnpoint()
	homeDeviceServiceImpl := hDService.NewHomeDeviceServiceImpl2(dao.HomeDeviceDaoImpl{DynamoDbApi: svc})

	_, err := HandleRequest(context.Background(), request, homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("Unexpected error running TestCreateHomeDevice_DeviceAlreadyExist. Error: %v", err.Error())
	}

	response2, err2 := HandleRequest(context.Background(), request, homeDeviceServiceImpl)

	if err2 != nil {
		t.Fatalf("Unexpected error running TestCreateHomeDevice_DeviceAlreadyExist. Error: %v", err.Error())
	}

	assert.Contains(t, response2.Body, "Device Already Exist")
	assert.Equal(t, response2.StatusCode, 400)
}
