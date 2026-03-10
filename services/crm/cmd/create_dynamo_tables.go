package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
	internalConfig "github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/logger"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(createTablesCommand)
}

var createTablesCommand = &cobra.Command{
	Use:     "create-tables",
	Short:   "Create all required DynamoDB tables",
	GroupID: "make",
	Run: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()

		dynamoCfg := internalConfig.GetConfig().DynamodbForms
		if dynamoCfg.DynamoEndpoint == "" {
			logger.Log.Fatal("DynamoDB endpoint not configured (dynamodbForms.dynamoEndpoint)")
		}

		ctx := context.Background()

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(dynamoCfg.AwsRegion),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: dynamoCfg.DynamoEndpoint}, nil
				})),
			config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "",
					Source: "Mock credentials used above for local instance",
				},
			}),
		)
		if err != nil {
			logger.Log.Fatal("Unable to load AWS SDK config", zap.Error(err))
		}

		client := dynamodb.NewFromConfig(cfg)

		tables := []dynamodb.CreateTableInput{
			{
				TableName: aws.String("forms"),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("PK"), AttributeType: types.ScalarAttributeTypeS},
					{AttributeName: aws.String("SK"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("SK"), KeyType: types.KeyTypeRange},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
			{
				TableName: aws.String("form_submissions"),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("PK"), AttributeType: types.ScalarAttributeTypeS},
					{AttributeName: aws.String("SK"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("SK"), KeyType: types.KeyTypeRange},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
			{
				TableName: aws.String("tags"),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("PK"), AttributeType: types.ScalarAttributeTypeS},
					{AttributeName: aws.String("SK"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("SK"), KeyType: types.KeyTypeRange},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		}

		for _, input := range tables {
			tableName := aws.ToString(input.TableName)
			_, err := client.CreateTable(ctx, &input)
			if err != nil {
				var resourceInUse *types.ResourceInUseException
				if errors.As(err, &resourceInUse) {
					logger.Log.Info(fmt.Sprintf("Table already exists, skipping: %s", tableName))
					continue
				}
				logger.Log.Fatal("Failed to create table", zap.String("table", tableName), zap.Error(err))
			}
			logger.Log.Info("Created table", zap.String("table", tableName))
		}

		logger.Log.Info("All tables created successfully")
	},
}
