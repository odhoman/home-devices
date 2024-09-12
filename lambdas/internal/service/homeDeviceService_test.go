package service

import (
	"context"
	"testing"

	"github.com/odhoman/home-devices/internal/constants"
	hdError "github.com/odhoman/home-devices/internal/error"
	hdMock "github.com/odhoman/home-devices/internal/mock"
	"github.com/odhoman/home-devices/internal/request"
	hdREsponse "github.com/odhoman/home-devices/internal/response"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

// Mock del DAO

func TestCreateHomeDevice_DeviceExist(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.CreateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	// Caso 1: El dispositivo ya existe
	mockDao.On("IsDeviceExist", ctx, deviceRequest.MAC, deviceRequest.HomeID).Return(true, (*hdError.HomeDeviceError)(nil))

	_, err := service.CreateHomeDevice(ctx, deviceRequest)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceAlreadyExistsCode, err.ErrorCode)
}

func TestCreateHomeDevice_ErrorVerifyingIfDeviceExist(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.CreateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	// Caso 1: El dispositivo ya existe
	mockDao.On("IsDeviceExist", ctx, deviceRequest.MAC, deviceRequest.HomeID).Return(false, &hdError.HomeDeviceError{ErrorCode: constants.ErrDeviceAlreadyExistsCode})

	_, err := service.CreateHomeDevice(ctx, deviceRequest)
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrDeviceAlreadyExistsCode, err.ErrorCode)
}

func TestCreateHomeDevice_Success(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.CreateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	// Caso 1: El dispositivo no existe
	mockDao.On("IsDeviceExist", ctx, deviceRequest.MAC, deviceRequest.HomeID).Return(false, (*hdError.HomeDeviceError)(nil))
	// Asegurarse de que el segundo retorno es nil pero del tipo correcto
	mockDao.On("SaveHomeDevice", ctx, deviceRequest).Return(&hdREsponse.HomdeDeviceResponse{}, (*hdError.HomeDeviceError)(nil))

	resp, err := service.CreateHomeDevice(ctx, deviceRequest)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestCreateHomeDevice_ErrorSavingDevice(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.CreateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	// Caso 1: El dispositivo no existe
	mockDao.On("IsDeviceExist", ctx, deviceRequest.MAC, deviceRequest.HomeID).Return(false, (*hdError.HomeDeviceError)(nil))
	mockDao.On("SaveHomeDevice", ctx, deviceRequest).Return(&hdREsponse.HomdeDeviceResponse{}, &hdError.HomeDeviceError{ErrorCode: "save_error"})

	_, err := service.CreateHomeDevice(ctx, deviceRequest)
	assert.NotNil(t, err)
	assert.Equal(t, "save_error", err.ErrorCode)
}

func TestUpdateHomeDevice_Success(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.UpdateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	mockDao.On("UpdateHomeDevice", ctx, deviceRequest, "id").Return((*hdError.HomeDeviceError)(nil))

	err := service.UpdateHomeDevice(ctx, deviceRequest, "id")
	assert.Nil(t, err)
}

func TestUpdateHomeDevice_NoFieldToUpdate(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.UpdateDeviceRequest{}

	err := service.UpdateHomeDevice(ctx, deviceRequest, "id")
	assert.NotNil(t, err)
	assert.Equal(t, constants.ErrNoFieldToUpdateCode, err.ErrorCode)
}

func TestUpdateHomeDevice_Error(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	deviceRequest := request.UpdateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	mockDao.On("UpdateHomeDevice", ctx, deviceRequest, "id").Return(&hdError.HomeDeviceError{ErrorCode: "save_error"})

	err := service.UpdateHomeDevice(ctx, deviceRequest, "id")
	assert.NotNil(t, err)
	assert.Equal(t, "save_error", err.ErrorCode)
}

func TestGetHomeDevice_Success(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	//deviceRequest := request.UpdateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	mockDao.On("GetHomeDevice", ctx, "id").Return(&hdREsponse.HomdeDeviceResponse{MAC: "00:11:22:33:44:55", HomeID: "home1"}, (*hdError.HomeDeviceError)(nil))

	response, err := service.GetHomeDevice(ctx, "id")

	assert.Nil(t, err)
	assert.Equal(t, "00:11:22:33:44:55", response.MAC)

}

func TestGetHomeDevice_Error(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()

	mockDao.On("GetHomeDevice", ctx, mock.Anything).Return(nil, &hdError.HomeDeviceError{ErrorCode: "get_error"})

	_, err := service.GetHomeDevice(ctx, "id")

	// Verificamos que se produjo un error
	assert.NotNil(t, err)
	// Comprobamos que el código de error es el esperado
	assert.Equal(t, "get_error", err.ErrorCode)
}

func TestDeleteHomeDevice_Success(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()
	//deviceRequest := request.UpdateDeviceRequest{MAC: "00:11:22:33:44:55", HomeID: "home1"}

	mockDao.On("DeleteHomeDevice", ctx, "id").Return((*hdError.HomeDeviceError)(nil))

	err := service.DeleteHomeDevice(ctx, "id")

	assert.Nil(t, err)
	//assert.Equal(t, "00:11:22:33:44:55", response.MAC)

}

func TestDeleteHomeDevice_Error(t *testing.T) {
	mockDao := new(hdMock.MockHomeDeviceDao)
	service := HomeDeviceServiceImpl{homeDeviceDao: mockDao}

	ctx := context.Background()

	mockDao.On("DeleteHomeDevice", ctx, mock.Anything).Return(&hdError.HomeDeviceError{ErrorCode: "delete_error"})

	err := service.DeleteHomeDevice(ctx, "id")

	// Verificamos que se produjo un error
	assert.NotNil(t, err)
	// Comprobamos que el código de error es el esperado
	assert.Equal(t, "delete_error", err.ErrorCode)
}
