package main

import (
	"context"
	"log"

	hDConstants "lambdas/common/constants"
	hDResponse "lambdas/common/response"
	hDService "lambdas/common/service"
	hDValidation "lambdas/common/validation"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func handleRequest(ctx context.Context, id string, deviceService hDService.HomeDeviceService) (events.APIGatewayProxyResponse, error) {

	if err := hDValidation.CheckEmptyString("id", id); err != nil {
		return hDResponse.BadRequestErrorAPIGatewayProxyResponseSingleMessage(err.Error()), nil
	}

	if err := deviceService.DeleteHomeDevice(ctx, id); err != nil {
		return getErrorResponse(err.ErrorCode), nil
	}

	return hDResponse.ReturnOKWithMessageAPIGatewayProxyResponse(200, "Device deleted"), nil
}

func getErrorResponse(errorCode string) events.APIGatewayProxyResponse {
	switch errorCode {
	case hDConstants.ErrDeviceNotFoundCode:
		return hDResponse.ReturnNotFoundErrorAPIGatewayProxyResponseSingleMessage("Device Not Found")
	default:
		return hDResponse.InternalServerErrorAPIGatewayProxyResponseSingleMessage("Internal Server error deleting a device")
	}

}

func main() {

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load SDK config for deleteDevice lambda function, %v", err)
		}

		return handleRequest(ctx, request.PathParameters["id"], hDService.NewHomeDeviceServiceImpl(cfg))
	})
}
