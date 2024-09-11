package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	hDRequest "lambdas/common/request"
	hDResponse "lambdas/common/response"
	hDService "lambdas/common/service"
	hDValidation "lambdas/common/validation"

	hDConstants "github.com/odhoman/home-devices/internal/constants"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func handleRequest(ctx context.Context, device hDRequest.CreateDeviceRequest, deviceService hDService.HomeDeviceService) (events.APIGatewayProxyResponse, error) {

	if valdationOutput := hDValidation.ValidateDeviceRequestStruct(device); len(valdationOutput) > 0 {
		return hDResponse.ReturnBadRequestErrorAPIGatewayProxyResponse(valdationOutput), nil
	}

	deviceCreated, err := deviceService.CreateHomeDevice(ctx, device)

	if err != nil {
		return getErrorResponse(err.ErrorCode), nil
	}

	return hDResponse.ReturnAPIGatewayProxyResponse(201, deviceCreated), nil
}

func main() {

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load SDK config for createDevice lambda function, %v", err)
		}

		var createDeviceRequest hDRequest.CreateDeviceRequest
		if err := json.Unmarshal([]byte(request.Body), &createDeviceRequest); err != nil {
			log.Printf("Error deserializing JSON for createDevice lambda function: %v", err)
			return hDResponse.BadRequestErrorAPIGatewayProxyResponseSingleMessage(fmt.Sprintf("Invalid request body: %v", err)), nil
		}

		return handleRequest(ctx, createDeviceRequest, hDService.NewHomeDeviceServiceImpl(cfg))
	})
}

func getErrorResponse(errorCode string) events.APIGatewayProxyResponse {
	switch errorCode {
	case hDConstants.ErrDeviceAlreadyExistsCode:
		return hDResponse.BadRequestErrorAPIGatewayProxyResponseSingleMessage("Device Already Exist")
	default:
		return hDResponse.InternalServerErrorAPIGatewayProxyResponseSingleMessage("Internal Server error creating a new device")
	}
}
