package main

import (
	"context"
	"encoding/json"
	"testing"

	hDConstants "github.com/odhoman/home-devices/internal/constants"
	hDError "github.com/odhoman/home-devices/internal/error"
	hDMock "github.com/odhoman/home-devices/internal/mock"
	hDResponse "github.com/odhoman/home-devices/internal/response"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRequest_Success(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	responseMessage := &hDResponse.MessageResponse{Message: "Device deleted"}

	id := uuid.New().String()

	mockService.On("DeleteHomeDevice", mock.Anything, id).Return(nil)

	response, err := HandleRequest(context.TODO(), id, mockService)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	expectedBody, _ := json.Marshal(responseMessage)
	assert.JSONEq(t, string(expectedBody), response.Body)

	mockService.AssertExpectations(t)
}

func TestHandleRequest_ValidationError(t *testing.T) {

	response, _ := HandleRequest(context.TODO(), "", new(hDMock.MockHomeDeviceService))

	assert.Equal(t, 400, response.StatusCode)

	assert.Contains(t, response.Body, "Field 'id' is empty. Please enter a value")

}

func TestHandleRequest_DeviceNotFound(t *testing.T) {

	id := uuid.New().String()

	mockService := new(hDMock.MockHomeDeviceService)
	mockService.On("DeleteHomeDevice", mock.Anything, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrDeviceNotFoundCode,
	})

	response, _ := HandleRequest(context.TODO(), id, mockService)

	assert.Equal(t, 404, response.StatusCode)
	assert.Contains(t, response.Body, "Device Not Found")
}

func TestHandleRequest_InternalServerError(t *testing.T) {

	id := uuid.New().String()

	mockService := new(hDMock.MockHomeDeviceService)
	mockService.On("DeleteHomeDevice", mock.Anything, id).Return(&hDError.HomeDeviceError{
		ErrorCode: hDConstants.ErrDeletingDeviceCode,
	})

	response, _ := HandleRequest(context.TODO(), id, mockService)

	assert.Equal(t, 500, response.StatusCode)
	assert.Contains(t, response.Body, "Internal Server error deleting a device")
}
