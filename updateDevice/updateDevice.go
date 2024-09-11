package main

import (
	"context"
	"encoding/json"
	"fmt"

	hDConstants "lambdas/common/constants"
	hDRequest "lambdas/common/request"
	hDResponse "lambdas/common/response"
	hDService "lambdas/common/service"
	hDValidation "lambdas/common/validation"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func handleRequest(ctx context.Context, device hDRequest.UpdateDeviceRequest, id string, deviceService hDService.HomeDeviceService) (events.APIGatewayProxyResponse, error) {

	if valdationOutput := hDValidation.ValidateDeviceRequestStruct(device); len(valdationOutput) > 0 {
		return hDResponse.ReturnBadRequestErrorAPIGatewayProxyResponse(valdationOutput), nil
	}

	if err := deviceService.UpdateHomeDevice(ctx, device, id); err != nil {
		return getErrorResponse(err.ErrorCode), nil
	}

	return hDResponse.ReturnOKWithMessageAPIGatewayProxyResponse(201, "Device updated"), nil
}

func getErrorResponse(errorCode string) events.APIGatewayProxyResponse {
	switch errorCode {
	case hDConstants.ErrDeviceNotFoundCode:
		return hDResponse.ReturnNotFoundErrorAPIGatewayProxyResponseSingleMessage("Device Not Found")
	case hDConstants.ErrNoFieldToUpdateCode:
		return hDResponse.BadRequestErrorAPIGatewayProxyResponseSingleMessage("Please enter a value property to update")
	default:
		return hDResponse.InternalServerErrorAPIGatewayProxyResponseSingleMessage("Internal Server error updating a device")
	}

}

func main() {

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load SDK config for homeDeviceListener lambda function, %v", err)
		}

		var updateDeviceRequest hDRequest.UpdateDeviceRequest
		if err := json.Unmarshal([]byte(request.Body), &updateDeviceRequest); err != nil {
			log.Printf("Error deserializing JSON: %v", err)
			return hDResponse.BadRequestErrorAPIGatewayProxyResponseSingleMessage(fmt.Sprintf("Invalid request body: %v", err)), nil
		}

		return handleRequest(ctx, updateDeviceRequest, request.PathParameters["id"], hDService.NewHomeDeviceServiceImpl(cfg))
	})
}
