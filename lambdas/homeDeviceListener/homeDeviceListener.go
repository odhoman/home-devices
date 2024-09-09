package main

import (
	"context"
	"encoding/json"
	"log"

	hDRequest "lambdas/common/request"
	hDService "lambdas/common/service"
	hDValidation "lambdas/common/validation"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type UpdateDeviceSQSMessage struct {
	ID     string `json:"id" validate:"required"`
	HomeID string `json:"homeId" validate:"required,min=5,max=30"`
}

type Device struct {
	ID        string `json:"id"`
	MAC       string `json:"mac"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	HomeID    string `json:"homeId"`
	CreatedAt int64  `json:"createdAt"`
}

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent, deviceService hDService.HomeDeviceService) {

	for _, message := range sqsEvent.Records {

		var updateDeviceSQSMessage UpdateDeviceSQSMessage
		if err := buildUpdateDeviceSQSMessage(message.Body, &updateDeviceSQSMessage); err != nil {
			log.Printf("Error parsing SQS message: %v", err)
			continue
		}

		if valdationOutput := hDValidation.ValidateAndResponseBadRequestErrors(updateDeviceSQSMessage); len(valdationOutput) > 0 {
			log.Printf("Validation Errors found in the request message : %v", valdationOutput)
			continue
		}

		deviceId := updateDeviceSQSMessage.ID
		homeId := updateDeviceSQSMessage.HomeID

		if err := deviceService.UpdateHomeDevice(ctx, hDRequest.UpdateDeviceRequest{
			HomeID: homeId,
		}, deviceId); err != nil {
			log.Printf("An error occurred updating a device for id %v - homeId %v: %v", deviceId, homeId, err.ErrorMessage)
			continue
		}

		log.Print("Invoking UpatedDevice done...")
	}
}

func buildUpdateDeviceSQSMessage(message string, updateDeviceSQSMessage *UpdateDeviceSQSMessage) error {
	return json.Unmarshal([]byte(message), &updateDeviceSQSMessage)
}

func main() {

	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load SDK config for homeDeviceListener lambda function, %v", err)
		}

		handleRequest(ctx, sqsEvent, hDService.NewHomeDeviceServiceImpl(cfg))
	})
}
