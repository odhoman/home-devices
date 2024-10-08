package dao

import (
	"context"
	"errors"
	"fmt"

	hdError "github.com/odhoman/home-devices/internal/error"
	utils "github.com/odhoman/home-devices/internal/utils"

	constants "github.com/odhoman/home-devices/internal/constants"
	request "github.com/odhoman/home-devices/internal/request"
	response "github.com/odhoman/home-devices/internal/response"

	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type HomeDeviceDao interface {
	IsDeviceExist(ctx context.Context, mac string, homeId string) (bool, *hdError.HomeDeviceError)
	SaveHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError)
	GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError)
	UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError
	DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError
}

type HomeDeviceDaoImpl struct {
	DynamoDbApi dynamoDbApi
}

func (hDDI HomeDeviceDaoImpl) IsDeviceExist(ctx context.Context, mac string, homeId string) (bool, *hdError.HomeDeviceError) {

	tableName, error := getValuePropertyOrError(constants.TableNameHomeDevicesProperty)
	if error != nil {
		return false, error
	}

	macHomeIdIndexName, error := getValuePropertyOrError(constants.MacHomeIdIndexNameProperty)
	if error != nil {
		return false, error
	}

	input := &dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              &macHomeIdIndexName,
		KeyConditionExpression: aws.String("mac = :mac and homeId = :homeId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":mac":    &types.AttributeValueMemberS{Value: mac},
			":homeId": &types.AttributeValueMemberS{Value: homeId},
		},
	}

	result, err := hDDI.DynamoDbApi.Query(ctx, input)

	if err != nil {
		fmt.Printf("Error querying the GSI: %v", err)
		return false, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrGettingDeviceCode,
			ErrorMessage: constants.ErrGettingDeviceMessage,
		}
	}

	return len(result.Items) > 0, nil

}

func (hDDI HomeDeviceDaoImpl) SaveHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {

	tableName, error := getValuePropertyOrError(constants.TableNameHomeDevicesProperty)
	if error != nil {
		return nil, error
	}

	id := uuid.New().String()
	now := time.Now().Unix()

	if _, err := hDDI.DynamoDbApi.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberS{Value: id},
			"mac":        &types.AttributeValueMemberS{Value: device.MAC},
			"name":       &types.AttributeValueMemberS{Value: device.Name},
			"type":       &types.AttributeValueMemberS{Value: device.Type},
			"homeId":     &types.AttributeValueMemberS{Value: device.HomeID},
			"createdAt":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", now)},
			"modifiedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", now)},
		},
	}); err != nil {
		fmt.Printf("Error putting item into DynamoDB: %v", err)
		return nil, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrDeviceNotCreatedErrorCode,
			ErrorMessage: constants.ErrDeviceNotCreatedErrorMessage,
		}
	}

	return &response.HomdeDeviceResponse{
		ID:         id,
		MAC:        device.MAC,
		Name:       device.Name,
		Type:       device.Type,
		HomeID:     device.HomeID,
		CreatedAt:  now,
		ModifiedAt: now,
	}, nil
}

func (hDDI HomeDeviceDaoImpl) GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {

	tableName, error := getValuePropertyOrError(constants.TableNameHomeDevicesProperty)
	if error != nil {
		return nil, error
	}

	result, err := hDDI.DynamoDbApi.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		log.Printf("Error getting item with it %v DynamoDB: %v", id, err)
		return nil, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrGettingDeviceCode,
			ErrorMessage: constants.ErrDeviceNotCreatedErrorMessage,
		}
	}

	if result.Item == nil {
		return nil, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrDeviceNotFoundCode,
			ErrorMessage: constants.ErrDeviceNotFoundMessage,
		}
	}

	device := mapDynamoDBItemToDeviceResponse(result.Item)

	return &device, nil
}

func (hDDI HomeDeviceDaoImpl) UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError {

	updateInput, error := buidUpdateInput(device, id)
	if error != nil {
		return error
	}

	if _, err := hDDI.DynamoDbApi.UpdateItem(ctx, updateInput); err != nil {

		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return &hdError.HomeDeviceError{
				ErrorCode:    constants.ErrDeviceNotFoundCode,
				ErrorMessage: constants.ErrDeviceNotFoundMessage,
			}
		}

		log.Printf("Error updating item with id %v into DynamoDB: %v", id, err)
		return &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrUpdatingDeviceCode,
			ErrorMessage: constants.ErrUpdatingDeviceMessage,
		}
	}

	return nil
}

func (hDDI HomeDeviceDaoImpl) DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError {

	tableName, error := getValuePropertyOrError(constants.TableNameHomeDevicesProperty)
	if error != nil {
		return error
	}

	if _, err := hDDI.DynamoDbApi.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	}); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			log.Printf("Record with id %v does not exist, delete failed", id)
			return &hdError.HomeDeviceError{
				ErrorCode:    constants.ErrDeviceNotFoundCode,
				ErrorMessage: constants.ErrDeviceNotFoundMessage,
			}
		}

		log.Printf("Error deleting item with id %v into DynamoDB: %v", id, err)
		return &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrDeletingDeviceCode,
			ErrorMessage: constants.ErrDeletingDeviceMessage,
		}

	}

	return nil
}

func buidUpdateInput(device request.UpdateDeviceRequest, id string) (*dynamodb.UpdateItemInput, *hdError.HomeDeviceError) {

	tableName, error := getValuePropertyOrError(constants.TableNameHomeDevicesProperty)
	if error != nil {
		return nil, error
	}

	updateExpression := "SET modifiedAt = :modifiedAt"
	expressionAttributeValues := map[string]types.AttributeValue{
		":modifiedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
	}

	expressionAttributeNames := map[string]string{}

	if device.MAC != "" {
		updateExpression += ", mac = :mac"
		expressionAttributeValues[":mac"] = &types.AttributeValueMemberS{Value: device.MAC}
	}

	if device.Name != "" {
		updateExpression += ", #name = :name"
		expressionAttributeValues[":name"] = &types.AttributeValueMemberS{Value: device.Name}
		expressionAttributeNames["#name"] = "name"
	}

	if device.Type != "" {
		updateExpression += ", #type = :type"
		expressionAttributeValues[":type"] = &types.AttributeValueMemberS{Value: device.Type}
		expressionAttributeNames["#type"] = "type"
	}

	if device.HomeID != "" {
		updateExpression += ", homeId = :homeId"
		expressionAttributeValues[":homeId"] = &types.AttributeValueMemberS{Value: device.HomeID}
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ConditionExpression:       aws.String("attribute_exists(id)"),
	}

	if len(expressionAttributeNames) > 0 {
		updateInput.ExpressionAttributeNames = expressionAttributeNames
	}

	return updateInput, nil
}

func getValuePropertyOrError(fieldName string) (string, *hdError.HomeDeviceError) {
	value, error := utils.GetValueProperty(fieldName)

	if error != nil {
		return "", &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrGettingConfigCode,
			ErrorMessage: constants.ErrGettingConfigMessage,
		}
	}

	return value, nil
}

func mapDynamoDBItemToDeviceResponse(item map[string]types.AttributeValue) response.HomdeDeviceResponse {
	return response.HomdeDeviceResponse{
		ID:         getStringAttribute(item, "id"),
		MAC:        getStringAttribute(item, "mac"),
		Name:       getStringAttribute(item, "name"),
		Type:       getStringAttribute(item, "type"),
		HomeID:     getStringAttribute(item, "homeId"),
		CreatedAt:  getInt64Attribute(item, "createdAt"),
		ModifiedAt: getInt64Attribute(item, "modifiedAt"),
	}
}

func getStringAttribute(item map[string]types.AttributeValue, key string) string {
	if v, ok := item[key].(*types.AttributeValueMemberS); ok {
		return v.Value
	}
	return ""
}

func getInt64Attribute(item map[string]types.AttributeValue, key string) int64 {
	if v, ok := item[key].(*types.AttributeValueMemberN); ok {
		if value, err := strconv.ParseInt(v.Value, 10, 64); err == nil {
			return value
		}
	}
	return 0
}
