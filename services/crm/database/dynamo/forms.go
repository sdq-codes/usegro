package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/usegro/services/crm/internal/logger"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DynamoClient *dynamodb.Client
var m sync.Mutex

// InitDynamoFormClient initializes a global DynamoDB client (singleton-style).
func InitDynamoFormClient(endpoint, region string) error {
	m.Lock()
	defer m.Unlock()

	if DynamoClient != nil {
		return nil // already initialized
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
		return err
	}

	DynamoClient = dynamodb.NewFromConfig(cfg)
	return nil
}

// GetDynamoClient returns the initialized DynamoDB client.
func GetDynamoClient() *dynamodb.Client {
	if DynamoClient == nil {
		panic("DynamoDB client not initialized. Call InitDynamoClient first.")
	}
	return DynamoClient
}

// CloseDynamoClient cleans up resources (not strictly needed, but mirrors pgx pattern).
func CloseDynamoClient() {
	m.Lock()
	defer m.Unlock()

	if DynamoClient != nil {
		DynamoClient = nil
	}
}

func CreateTable() {
	logger.Log.Info("Creating table form")
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("forms"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := DynamoClient.CreateTable(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}
