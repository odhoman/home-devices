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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRequest_Success(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	request := hDRequest.CreateDeviceRequest{
		MAC:    "00-1A-2B-3C-4D-5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	device := &hDResponse.HomdeDeviceResponse{
		ID:         uuid.New().String(),
		MAC:        request.MAC,
		Name:       request.Name,
		Type:       request.Type,
		HomeID:     request.HomeID,
		CreatedAt:  time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	}

	mockService.On("CreateHomeDevice", mock.Anything, request).Return(device, nil)

	response, err := handleRequest(context.TODO(), request, mockService)

	assert.NoError(t, err)
	assert.Equal(t, 201, response.StatusCode)

	expectedBody, _ := json.Marshal(device)
	assert.JSONEq(t, string(expectedBody), response.Body)

	mockService.AssertExpectations(t)
}

func TestHandleRequest_ValidationErrorRequiredFields(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		HomeID: "home123",
		MAC:    "WRONG MAC ADDRESS",
	}

	response, _ := handleRequest(context.TODO(), request, new(hDMock.MockHomeDeviceService))

	assert.Equal(t, 400, response.StatusCode)

	assert.Contains(t, response.Body, "Please enter a valid MAC address")
	assert.Contains(t, response.Body, "Validation failed for field 'Type': required")
	assert.Contains(t, response.Body, "Validation failed for field 'Name': required")
}

func TestHandleRequest_ValidationErrorFieldsLength(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "mac super large mac super large  mac super large mac super large mac super large mac super large mac super large ",
		Name:   "name super large name super large name super large name super large name super large name super large name super large ",
		Type:   "type super large type super large type super large type super large type super large type super large type super large ",
		HomeID: "homeId super large homeId super large homeId super large homeId super large homeId super large homeId super large homeId super large ",
	}

	response, _ := handleRequest(context.TODO(), request, new(hDMock.MockHomeDeviceService))

	assert.Equal(t, 400, response.StatusCode)

	assert.Contains(t, response.Body, "MAC address must be between 12 and 17 characters")
	assert.Contains(t, response.Body, "Name must be between 3 and 50 characters")
	assert.Contains(t, response.Body, "Type must be between 3 and 20 characters")
	assert.Contains(t, response.Body, "Home ID must be between 5 and 30 characters")
}

func TestHandleRequest_AlreadyExistError(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "00:1B:44:11:3A:B7",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	mockService := new(hDMock.MockHomeDeviceService)
	mockService.On("CreateHomeDevice", mock.Anything, request).Return(nil, &hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrDeviceAlreadyExistsCode,
	})

	response, _ := handleRequest(context.TODO(), request, mockService)

	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Device Already Exist")
}

func TestHandleRequest_InternalServerError(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "00:1B:44:11:3A:B7",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	mockService := new(hDMock.MockHomeDeviceService)
	mockService.On("CreateHomeDevice", mock.Anything, request).Return(nil, &hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrDeviceNotCreatedErrorCode,
	})

	response, _ := handleRequest(context.TODO(), request, mockService)

	assert.Equal(t, 500, response.StatusCode)
	assert.Contains(t, response.Body, "Internal Server error creating a new device")
}
