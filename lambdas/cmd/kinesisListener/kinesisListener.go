package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, kinesisEvent events.KinesisEvent) (string, error) {
	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis

		// Los datos ya están en formato []byte, no necesitas decodificarlos de base64
		data := kinesisRecord.Data

		// Loguear la Partition Key y los datos recibidos
		log.Printf("Partition Key: %s", kinesisRecord.PartitionKey)
		log.Printf("Data: %s", string(data))
	}

	return "Processed Kinesis Event", nil
}

func main() {
	lambda.Start(HandleRequest)
}
