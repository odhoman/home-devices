package service

import (
	"context"

	hdError "github.com/odhoman/home-devices/internal/error"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	constants "github.com/odhoman/home-devices/internal/constants"
	dao "github.com/odhoman/home-devices/internal/dao"
	request "github.com/odhoman/home-devices/internal/request"
	response "github.com/odhoman/home-devices/internal/response"
)

type HomeDeviceService interface {
	CreateHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError)
	GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError)
	UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError
	DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError
}

type HomeDeviceServiceImpl struct {
	homeDeviceDao dao.HomeDeviceDao
}

func (hDDI HomeDeviceServiceImpl) CreateHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {

	dao := hDDI.homeDeviceDao

	isExist, err := dao.IsDeviceExist(ctx, device.MAC, device.HomeID)
	if err != nil {
		return nil, err
	}

	if isExist {
		return nil, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrDeviceAlreadyExistsCode,
			ErrorMessage: constants.ErrDeviceAlreadyExistsMessage,
		}
	}

	response, saveDeviceError := dao.SaveHomeDevice(ctx, device)
	if saveDeviceError != nil {
		return nil, saveDeviceError
	}

	return response, nil

}

func (hDDI HomeDeviceServiceImpl) GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {

	dao := hDDI.homeDeviceDao

	result, err := dao.GetHomeDevice(ctx, id)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (hDDI HomeDeviceServiceImpl) UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError {

	if device.MAC == "" && device.Name == "" && device.Type == "" && device.HomeID == "" {
		return &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrNoFieldToUpdateCode,
			ErrorMessage: constants.ErrNoFieldToUpdateMessage,
		}
	}

	dao := hDDI.homeDeviceDao

	return dao.UpdateHomeDevice(ctx, device, id)
}

func (hDDI HomeDeviceServiceImpl) DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError {
	dao := hDDI.homeDeviceDao
	return dao.DeleteHomeDevice(ctx, id)
}

func NewHomeDeviceServiceImplFromConfig2(cfg aws.Config) HomeDeviceService {
	client := dynamodb.NewFromConfig(cfg)
	dao := dao.HomeDeviceDaoImpl{DynamoDbApi: client}
	return NewHomeDeviceServiceImpl2(dao)
}

func NewHomeDeviceServiceImpl2(dao dao.HomeDeviceDao) HomeDeviceService {
	return HomeDeviceServiceImpl{homeDeviceDao: dao}
}
