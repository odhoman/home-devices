package common

import (
	"context"
	"errors"
	"fmt"
	constants "lambdas/common/constants"
	hdError "lambdas/common/error"
	request "lambdas/common/request"
	response "lambdas/common/response"
	utils "lambdas/common/utils"

	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type HomeDeviceService interface {
	CreateHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError)
	GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError)
	UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError
	DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError
}

type HomeDeviceServiceImpl struct {
	DynamoDbApi dynamoDbApi
}

func (hDDI HomeDeviceServiceImpl) CreateHomeDevice(ctx context.Context, device request.CreateDeviceRequest) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {

	tableName, error := getValuePropertyOrError(constants.TableNameHomeDevicesProperty)
	if error != nil {
		return nil, error
	}

	macHomeIdIndexName, error := getValuePropertyOrError(constants.MacHomeIdIndexNameProperty)
	if error != nil {
		return nil, error
	}

	input := &dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              &macHomeIdIndexName,
		KeyConditionExpression: aws.String("mac = :mac and homeId = :homeId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":mac":    &types.AttributeValueMemberS{Value: device.MAC},
			":homeId": &types.AttributeValueMemberS{Value: device.HomeID},
		},
	}

	result, err := hDDI.DynamoDbApi.Query(ctx, input)
	if err != nil {
		log.Printf("Error querying the GSI: %v", err)
		return nil, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrGettingDeviceCode,
			ErrorMessage: constants.ErrGettingDeviceMessage,
		}
	}

	if len(result.Items) > 0 {
		return nil, &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrDeviceAlreadyExistsCode,
			ErrorMessage: constants.ErrDeviceAlreadyExistsMessage,
		}
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
		log.Printf("Error putting item into DynamoDB: %v", err)
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

func (hDDI HomeDeviceServiceImpl) GetHomeDevice(ctx context.Context, id string) (*response.HomdeDeviceResponse, *hdError.HomeDeviceError) {

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

	device := response.HomdeDeviceResponse{}

	if v, ok := result.Item["id"].(*types.AttributeValueMemberS); ok {
		device.ID = v.Value
	}
	if v, ok := result.Item["mac"].(*types.AttributeValueMemberS); ok {
		device.MAC = v.Value
	}
	if v, ok := result.Item["name"].(*types.AttributeValueMemberS); ok {
		device.Name = v.Value
	}
	if v, ok := result.Item["type"].(*types.AttributeValueMemberS); ok {
		device.Type = v.Value
	}
	if v, ok := result.Item["homeId"].(*types.AttributeValueMemberS); ok {
		device.HomeID = v.Value
	}
	if v, ok := result.Item["createdAt"].(*types.AttributeValueMemberN); ok {
		device.CreatedAt, _ = strconv.ParseInt(v.Value, 10, 64)
	}
	if v, ok := result.Item["modifiedAt"].(*types.AttributeValueMemberN); ok {
		device.ModifiedAt, _ = strconv.ParseInt(v.Value, 10, 64)
	}

	return &device, nil
}

func (hDDI HomeDeviceServiceImpl) UpdateHomeDevice(ctx context.Context, device request.UpdateDeviceRequest, id string) *hdError.HomeDeviceError {

	if device.MAC == "" && device.Name == "" && device.Type == "" && device.HomeID == "" {
		return &hdError.HomeDeviceError{
			ErrorCode:    constants.ErrNoFieldToUpdateCode,
			ErrorMessage: constants.ErrNoFieldToUpdateMessage,
		}
	}

	updateInput, error := buidUpdateInput(device, id)
	if error != nil {
		return error
	}

	if _, err := hDDI.DynamoDbApi.UpdateItem(ctx, updateInput); err != nil {

		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {

			log.Printf("Record with id %v does not exist, update failed", updateInput.Key["id"])
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

func (hDDI HomeDeviceServiceImpl) DeleteHomeDevice(ctx context.Context, id string) *hdError.HomeDeviceError {

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

		log.Printf("DeleteItem  " + err.Error())
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

	log.Printf("No hay error")
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

func NewHomeDeviceServiceImpl(cfg aws.Config) HomeDeviceService {
	return HomeDeviceServiceImpl{DynamoDbApi: dynamodb.NewFromConfig(cfg)}
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
