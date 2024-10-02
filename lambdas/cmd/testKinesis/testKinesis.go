package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func main() {
	// Load custom configuration for LocalStack
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // Region for LocalStack
		//config.WithCredentialsProvider(aws.NewStaticCredentialsProvider("fakeAccessKey", "fakeSecretKey", "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "http://localhost:4566", // Change to the URL where your LocalStack runs
				SigningRegion: "us-east-1",
			}, nil
		})),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create a Kinesis client pointing to LocalStack
	kinesisClient := kinesis.NewFromConfig(cfg)

	streamName := "HomeDevicesStack-HomeDevicesKinesisStream-dc38db91" // Replace with your stream in LocalStack

	partitionKeys := []string{"MiPartitionKey1", "MiPartitionKey2", "MiPartitionKey3", "MiPartitionKey4"}

	// Channel to capture interrupt signals (Ctrl + C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Starting to send data to Kinesis...")

	// Infinite loop
	for {
		select {
		case <-stop:
			fmt.Println("Interrupted by user, terminating the program...")
			return
		default:
			// Randomly choose a PartitionKey
			partitionKey := partitionKeys[rand.Intn(len(partitionKeys))]

			// Data to be sent
			data := fmt.Sprintf("Message in partition key: %s", partitionKey)

			// Send the record to Kinesis
			_, err := kinesisClient.PutRecord(context.TODO(), &kinesis.PutRecordInput{
				StreamName:   aws.String(streamName),
				Data:         []byte(data),
				PartitionKey: aws.String(partitionKey),
			})

			if err != nil {
				log.Printf("Error sending data to Kinesis: %v", err)
			} else {
				log.Printf("Data sent to partition key %s: %s", partitionKey, data)
			}

			//time.Sleep(200 * time.Millisecond)
		}
	}
}
