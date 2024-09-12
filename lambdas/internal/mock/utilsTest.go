package mock

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	hDConstants "github.com/odhoman/home-devices/internal/constants"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func MockEnvVars() {
	os.Setenv(hDConstants.TableNameHomeDevicesProperty, "table")
	os.Setenv(hDConstants.MacHomeIdIndexNameProperty, "index")
}

func ClearEnvVars() {
	os.Setenv(hDConstants.TableNameHomeDevicesProperty, "")
	os.Setenv(hDConstants.MacHomeIdIndexNameProperty, "")
}

func Setup(t *testing.M) (string, func(t *testing.M)) {

	endpoint := ""
	MockEnvVars()

	fmt.Println("Setup init...")

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog("Initializing DynamoDB Local with the following configuration"),
	}

	fmt.Println("Docker container from image amazon/dynamodb-local created")

	var err error
	dynamodbC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	host, err := dynamodbC.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get container host: %s", err)
	}
	port, err := dynamodbC.MappedPort(ctx, "8000")
	if err != nil {
		log.Fatalf("Could not get container port: %s", err)
	}

	// Crea la URL del endpoint de DynamoDB
	endpoint = fmt.Sprintf("http://%s:%s", host, port.Port())

	// Resolver de endpoint personalizado para DynamoDB Local
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		},
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithHTTPClient(&http.Client{
			Timeout: 60 * time.Second,
		}),
	)

	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	fmt.Println("Creting dynamo Table for testing")

	time.Sleep(2 * time.Second)

	_, err = svc.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String("HomeDevices"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("mac"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("homeId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("MacHomeIdIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("mac"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("homeId"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(100),
					WriteCapacityUnits: aws.Int64(100),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(100),
			WriteCapacityUnits: aws.Int64(100),
		},
	})

	if err != nil {
		log.Fatalf("Failed to create table, %v", err)
	}

	os.Setenv(hDConstants.TableNameHomeDevicesProperty, "HomeDevices")
	os.Setenv(hDConstants.MacHomeIdIndexNameProperty, "MacHomeIdIndex")

	fmt.Println("Setup finished...")

	return endpoint, func(m *testing.M) {
		ClearEnvVars()
		if err := dynamodbC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop DynamoDB: %s", err)
		}
	}

}

func GetDynamoConnectionTestFromEnpoint() *dynamodb.Client {
	ep := os.Getenv("dynamoDBLocalEnpoint")
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					URL: ep,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		},
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg)

}

func RunTestMain(m *testing.M) {
	endpointReturned, tearDown := Setup(m)
	os.Setenv("dynamoDBLocalEnpoint", endpointReturned)
	code := m.Run()
	tearDown(m)
	os.Exit(code)
}
