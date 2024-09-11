package common

import (
	"context"

	hdError "lambdas/common/error"
	request "lambdas/common/request"
	response "lambdas/common/response"

	"github.com/stretchr/testify/mock"
)

type MockHomeDeviceService struct {
	mock.Mock
}

func (m *MockHomeDeviceService) CreateHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {
	args := m.Called(ctx, device)
	if args.Get(0) != nil {
		return args.Get(0).(*response.HomdeDeviceResponse), nil
	}
	return nil, args.Get(1).(*hdError.HomeDeviceError)
}

func (m *MockHomeDeviceService) GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*response.HomdeDeviceResponse), nil
	}
	return nil, args.Get(1).(*hdError.HomeDeviceError)
}

func (m *MockHomeDeviceService) UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError {
	args := m.Called(ctx, device, id)
	if args.Get(0) != nil {
		return args.Get(0).(*hdError.HomeDeviceError)
	}
	return nil
}

func (m *MockHomeDeviceService) DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*hdError.HomeDeviceError)
	}
	return nil
}
