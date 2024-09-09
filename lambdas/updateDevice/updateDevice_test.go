package main

import (
	"context"
	"encoding/json"
	hDConstants "lambdas/common/constants"
	hDError "lambdas/common/error"
	hDMock "lambdas/common/mock"
	hDRequest "lambdas/common/request"
	hDResponse "lambdas/common/response"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRequest_Success(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "BB:WW:SS:AR:GA:UF",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	responseMessage := &hDResponse.MessageResponse{Message: "Device updated"}

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(nil)

	response, err := handleRequest(context.TODO(), request, id, mockService)

	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	expectedBody, _ := json.Marshal(responseMessage)
	assert.JSONEq(t, string(expectedBody), response.Body)

	mockService.AssertExpectations(t)
}

func TestHandleRequest_ValidationErrorFieldsLength(t *testing.T) {

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "mac super large mac super large  mac super large mac super large mac super large mac super large mac super large ",
		Name:   "name super large name super large name super large name super large name super large name super large name super large ",
		Type:   "type super large type super large type super large type super large type super large type super large type super large ",
		HomeID: "homeId super large homeId super large homeId super large homeId super large homeId super large homeId super large homeId super large ",
	}

	response, _ := handleRequest(context.TODO(), request, uuid.New().String(), new(hDMock.MockHomeDeviceService))

	assert.Equal(t, 400, response.StatusCode)

	assert.Contains(t, response.Body, "MAC address must be between 12 and 17 characters")
	assert.Contains(t, response.Body, "Name must be between 3 and 50 characters")
	assert.Contains(t, response.Body, "Type must be between 3 and 20 characters")
	assert.Contains(t, response.Body, "Home ID must be between 5 and 30 characters")
}

func TestHandleRequest_DeviceNotFound(t *testing.T) {

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "BB:WW:SS:AR:GA:UF",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	mockService := new(hDMock.MockHomeDeviceService)

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrDeviceNotFoundCode,
	})

	response, _ := handleRequest(context.TODO(), request, id, mockService)

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

	response, _ := handleRequest(context.TODO(), request, id, mockService)

	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Please enter a value property to update")
}

func TestHandleRequest_InternalServerError(t *testing.T) {

	id := uuid.New().String()

	request := hDRequest.UpdateDeviceRequest{
		MAC:    "BB:WW:SS:AR:GA:UF",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	mockService := new(hDMock.MockHomeDeviceService)

	mockService.On("UpdateHomeDevice", mock.Anything, request, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrUpdatingDeviceCode,
	})

	response, _ := handleRequest(context.TODO(), request, id, mockService)

	assert.Equal(t, 500, response.StatusCode)
	assert.Contains(t, response.Body, "Internal Server error updating a device")
}
