package main

import (
	"context"
	"log"

	hDConstants "github.com/odhoman/home-devices/internal/constants"
	hDResponse "github.com/odhoman/home-devices/internal/response"
	hDService "github.com/odhoman/home-devices/internal/service"
	hDValidation "github.com/odhoman/home-devices/internal/validation"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func HandleRequest(ctx context.Context, id string, deviceService hDService.HomeDeviceService) (events.APIGatewayProxyResponse, error) {

	if err := hDValidation.CheckEmptyString("id", id); err != nil {
		return hDResponse.BadRequestErrorAPIGatewayProxyResponseSingleMessage(err.Error()), nil
	}

	device, err := deviceService.GetHomeDevice(ctx, id)

	if err != nil {
		return getErrorResponse(err.ErrorCode), nil
	}

	return hDResponse.ReturnAPIGatewayProxyResponse(200, device), nil
}

func getErrorResponse(errorCode string) events.APIGatewayProxyResponse {
	switch errorCode {
	case hDConstants.ErrDeviceNotFoundCode:
		return hDResponse.ReturnNotFoundErrorAPIGatewayProxyResponseSingleMessage("Device Not Found")
	case hDConstants.ErrGettingDeviceMessage:
		return hDResponse.InternalServerErrorAPIGatewayProxyResponseSingleMessage("Internal Server error getting the device")
	default:
		return hDResponse.InternalServerErrorAPIGatewayProxyResponseSingleMessage("Internal Server error getting the device")
	}

}

func main() {
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load SDK config for getDevice lambda function, %v", err)
		}

		return HandleRequest(ctx, request.PathParameters["id"], hDService.NewHomeDeviceServiceImplFromConfig2(cfg))
	})
}
