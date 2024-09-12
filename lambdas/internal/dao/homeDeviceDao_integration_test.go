package dao

import (
	"context"
	"fmt"
	"testing"

	hDConstants "github.com/odhoman/home-devices/internal/constants"
	hdError "github.com/odhoman/home-devices/internal/error"
	"github.com/odhoman/home-devices/internal/mock"
	hDRequest "github.com/odhoman/home-devices/internal/request"
	hDResponse "github.com/odhoman/home-devices/internal/response"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mock.RunTestMain(m)
}

func TestIsDeviceExist_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "00:1A:2B:3C:4D:5E",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	ctx := context.Background()
	homeDeviceDaoImpl := createHomeDeviceDaoImpl()
	response, err := executeSaveHomeDevice(ctx, request, homeDeviceDaoImpl)

	if err != nil {
		fmt.Println(err.ErrorCode)
		return
	}

	IsDeviceExist, err := homeDeviceDaoImpl.IsDeviceExist(ctx, response.MAC, response.HomeID)

	if err != nil {
		t.Fatalf("expected a bool but got an error %v", err.ErrorCode)
	}

	assert.True(t, IsDeviceExist)
	assert.Equal(t, request.MAC, response.MAC)
}

func TestIsDeviceExist_False_Success(t *testing.T) {

	ctx := context.Background()
	homeDeviceServiceImpl := createHomeDeviceDaoImpl()

	IsDeviceExist, err := homeDeviceServiceImpl.IsDeviceExist(ctx, "AB:CD:EF:01:23:45", "home412")

	if err != nil {
		t.Fatalf("expected a bool but got an error %v", err.ErrorCode)
	}

	assert.False(t, IsDeviceExist)
}

func TestSaveHomeDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "DE:AD:BE:EF:CA",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	homeDeviceServiceImpl := createHomeDeviceDaoImpl()
	response, err := executeSaveHomeDevice(context.TODO(), request, homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("expected a new home device but got an error %v", err.ErrorCode)
	}

	assert.Equal(t, request.HomeID, response.HomeID)
	assert.Equal(t, request.MAC, response.MAC)
}

func TestGetHomeDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "01:23:45:67:89",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	ctx := context.Background()
	homeDeviceServiceImpl := createHomeDeviceDaoImpl()
	response, err := executeSaveHomeDevice(ctx, request, homeDeviceServiceImpl)

	if err != nil {
		fmt.Println(err.ErrorCode)
		return
	}

	getHomeDeviceResponse, err := homeDeviceServiceImpl.GetHomeDevice(ctx, response.ID)

	if err != nil {
		t.Fatalf("expected a home device but got an error %v", err.ErrorCode)
	}

	assert.Equal(t, getHomeDeviceResponse.HomeID, response.HomeID)
	assert.Equal(t, getHomeDeviceResponse.MAC, response.MAC)
	assert.Equal(t, getHomeDeviceResponse.Name, response.Name)
	assert.Equal(t, getHomeDeviceResponse.Type, response.Type)
	assert.Equal(t, getHomeDeviceResponse.ID, response.ID)

}

func TestGetHomeDevice_NoExist(t *testing.T) {

	homeDeviceServiceImpl := createHomeDeviceDaoImpl()

	_, err := homeDeviceServiceImpl.GetHomeDevice(context.Background(), "fakeId")

	if err == nil {
		t.Fatal("expected an error when getting a non-existent home device, but got nil")
	}

	assert.NotNil(t, err)
	assert.Equal(t, hDConstants.ErrDeviceNotFoundCode, err.ErrorCode)
}

func TestUpdateHomeDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "DE:AD:BE:EF:CA",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	updateRequest := hDRequest.UpdateDeviceRequest{
		MAC:    "99:AD:BE:EF:CA",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	ctx := context.TODO()
	homeDeviceServiceImpl := createHomeDeviceDaoImpl()
	executeSaveHomeDevice(ctx, request, homeDeviceServiceImpl)

	response, err := homeDeviceServiceImpl.SaveHomeDevice(context.Background(), request)

	if err != nil {
		t.Fatalf("expected a new home device, testing TestUpdateHomeDevice_Success but got an error %v", err.ErrorCode)
	}

	if err := homeDeviceServiceImpl.UpdateHomeDevice(context.Background(), updateRequest, response.ID); err != nil {
		t.Fatalf("expecting update a home device , testing TestUpdateHomeDevice_Success but got an error %v", err.ErrorCode)
	}

	assert.Nil(t, err)
}

func TestUpdateHomeDevice_NoExist(t *testing.T) {

	updateRequest := hDRequest.UpdateDeviceRequest{
		MAC:    "99:AD:BE:EF:CA",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}

	homeDeviceServiceImpl := createHomeDeviceDaoImpl()
	err := homeDeviceServiceImpl.UpdateHomeDevice(context.Background(), updateRequest, "fakeID")

	if err == nil {
		t.Fatal("expected an error when updating a non-existent home device, but got nil")
	}

	assert.Equal(t, hDConstants.ErrDeviceNotFoundCode, err.ErrorCode)
}

func TestDeleteHomeDevice_Success(t *testing.T) {

	request := hDRequest.CreateDeviceRequest{
		MAC:    "77:1B:2C:3D:4E:88",
		Name:   "Living Room Light",
		Type:   "light",
		HomeID: "home12122",
	}
	ctx := context.Background()
	homeDeviceServiceImpl := createHomeDeviceDaoImpl()

	response, err := executeSaveHomeDevice(ctx, request, homeDeviceServiceImpl)

	if err != nil {
		t.Fatalf("expected a nil error creating a new device to test TestDeleteHomeDevice_Success, but got %v", err.ErrorCode)
	}

	if err := homeDeviceServiceImpl.DeleteHomeDevice(context.Background(), response.ID); err != nil {
		t.Fatalf("expected a nil error when deleting a home device, but got %v", err.ErrorCode)
	}

}

func TestDeleteHomeDevice_NoExist(t *testing.T) {

	homeDeviceServiceImpl := createHomeDeviceDaoImpl()
	err := homeDeviceServiceImpl.DeleteHomeDevice(context.Background(), "fakeID")

	if err == nil {
		t.Fatal("expected an error when deleting a non-existent home device, but got nil")
	}

	assert.Equal(t, hDConstants.ErrDeviceNotFoundCode, err.ErrorCode)

}

func createHomeDeviceDaoImpl() HomeDeviceDaoImpl {
	svc := mock.GetDynamoConnectionTestFromEnpoint()
	return HomeDeviceDaoImpl{DynamoDbApi: svc}
}

func executeSaveHomeDevice(ctx context.Context, request hDRequest.CreateDeviceRequest, homeDeviceServiceImpl HomeDeviceDaoImpl) (*hDResponse.HomdeDeviceResponse, *hdError.HomeDeviceError) {
	return homeDeviceServiceImpl.SaveHomeDevice(ctx, request)
}
