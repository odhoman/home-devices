package mock

import (
	"context"

	hdError "github.com/odhoman/home-devices/internal/error"
	request "github.com/odhoman/home-devices/internal/request"
	hdREsponse "github.com/odhoman/home-devices/internal/response"
	"github.com/stretchr/testify/mock"
)

type MockHomeDeviceDao struct {
	mock.Mock
}

func (m *MockHomeDeviceDao) IsDeviceExist(ctx context.Context, MAC, HomeID string) (bool, *hdError.HomeDeviceError) {
	args := m.Called(ctx, MAC, HomeID)
	return args.Bool(0), args.Get(1).(*hdError.HomeDeviceError)
}

func (m *MockHomeDeviceDao) SaveHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*hdREsponse.HomdeDeviceResponse, *hdError.HomeDeviceError) {
	args := m.Called(ctx, device)
	return args.Get(0).(*hdREsponse.HomdeDeviceResponse), args.Get(1).(*hdError.HomeDeviceError)
}

func (m *MockHomeDeviceDao) DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*hdError.HomeDeviceError)
	}
	return nil
}

func (m *MockHomeDeviceDao) GetHomeDevice(ctx context.Context, id string) (*hdREsponse.HomdeDeviceResponse, *hdError.HomeDeviceError) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*hdREsponse.HomdeDeviceResponse), nil
	}
	return nil, args.Get(1).(*hdError.HomeDeviceError)
}

func (m *MockHomeDeviceDao) UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError {
	args := m.Called(ctx, device, id)
	if args.Get(0) != nil {
		return args.Get(0).(*hdError.HomeDeviceError)
	}
	return nil
}
