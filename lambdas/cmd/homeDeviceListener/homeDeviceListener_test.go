package main

import (
	"context"

	"testing"

	hDError "github.com/odhoman/home-devices/internal/error"
	hDMock "github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/mock"
)

func TestHandleRequest_Success(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	mockService.On("UpdateHomeDevice", mock.Anything, hDRequest.UpdateDeviceRequest{HomeID: "home12345"}, "device123").Return(nil)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"id":"device123", "homeId":"home12345"}`,
			},
		},
	}

	HandleRequest(context.TODO(), sqsEvent, mockService)

	mockService.AssertCalled(t, "UpdateHomeDevice", mock.Anything, hDRequest.UpdateDeviceRequest{HomeID: "home12345"}, "device123")
}

func TestHandleRequest_UnmarshalError(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{ invalid json }`,
			},
		},
	}

	HandleRequest(context.TODO(), sqsEvent, mockService)

	mockService.AssertNotCalled(t, "UpdateHomeDevice")
}

func TestHandleRequest_ValidationError(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"id":"", "homeId":""}`,
			},
		},
	}

	HandleRequest(context.TODO(), sqsEvent, mockService)

	mockService.AssertNotCalled(t, "UpdateHomeDevice")
}

func TestHandleRequest_UpdateDeviceError(t *testing.T) {
	mockService := new(hDMock.MockHomeDeviceService)

	updateError := &hDError.HomeDeviceError{
		ErrorCode:    "InternalError",
		ErrorMessage: "Unable to update device",
	}

	mockService.On("UpdateHomeDevice", mock.Anything, hDRequest.UpdateDeviceRequest{HomeID: "home12345"}, "device123").Return(updateError)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"id":"device123", "homeId":"home12345"}`,
			},
		},
	}

	HandleRequest(context.TODO(), sqsEvent, mockService)

	mockService.AssertCalled(t, "UpdateHomeDevice", mock.Anything, hDRequest.UpdateDeviceRequest{HomeID: "home12345"}, "device123")
}
