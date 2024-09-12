package main

import (
	"context"
	"encoding/json"
	"testing"

	hDConstants "github.com/odhoman/home-devices/internal/constants"
	hDError "github.com/odhoman/home-devices/internal/error"
	hDMock "github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"
	hDResponse "github.com/odhoman/home-devices/internal/response"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRequest_Success(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "00-1A-2B-3C-4D-5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	responseMessage := &hDResponse.MessageResponse{Message: "Device updated"}

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(nil)

	response, err := HandleRequest(context.TODO(), request, id, mockService)

	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	expectedBody, _ := json.Marshal(responseMessage)
	assert.JSONEq(t, string(expectedBody), response.Body)

	mockService.AssertExpectations(t)
}

func TestHandleRequest_ValidationErrorFieldsLength(t *testing.T) {

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "001A-2B-3C-4D-5E",
		Name:   "name super large name super large name super large name super large name super large name super large name super large ",
		Type:   "type super large type super large type super large type super large type super large type super large type super large ",
		HomeID: "homeId super large homeId super large homeId super large homeId super large homeId super large homeId super large homeId super large ",
	}

	response, _ := HandleRequest(context.TODO(), request, uuid.New().String(), new(hDMock.MockHomeDeviceService))

	assert.Equal(t, 400, response.StatusCode)

	assert.Contains(t, response.Body, "Please enter a valid MAC address")
	assert.Contains(t, response.Body, "Name must be between 3 and 50 characters")
	assert.Contains(t, response.Body, "Type must be between 3 and 20 characters")
	assert.Contains(t, response.Body, "Home ID must be between 5 and 30 characters")
}

func TestHandleRequest_DeviceNotFound(t *testing.T) {

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "00-1A-2B-3C-4D-5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	mockService := new(hDMock.MockHomeDeviceService)

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrDeviceNotFoundCode,
	})

	response, _ := HandleRequest(context.TODO(), request, id, mockService)

	assert.Equal(t, 404, response.StatusCode)
	assert.Contains(t, response.Body, "Device Not Found")
}

func TestHandleRequest_NoFiledToUpdate(t *testing.T) {

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{}

	mockService := new(hDMock.MockHomeDeviceService)

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrNoFieldToUpdateCode,
	})

	response, _ := HandleRequest(context.TODO(), request, id, mockService)

	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Please enter a value property to update")
}

func TestHandleRequest_InternalServerError(t *testing.T) {

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "00-1A-2B-3C-4D-5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	mockService := new(hDMock.MockHomeDeviceService)

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrUpdatingDeviceCode,
	})

	response, _ := HandleRequest(context.TODO(), request, id, mockService)

	assert.Equal(t, 500, response.StatusCode)
	assert.Contains(t, response.Body, "Internal Server error updating a device")
}
