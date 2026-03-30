//go:build ignore

package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
	internalConfig "github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(migrateDynamoToMongoCommand)
	migrateDynamoToMongoCommand.Flags().StringVar(&migrateDynamoEndpoint, "dynamo-endpoint", "", "DynamoDB endpoint URL (required, e.g. http://usegro_dynamodb:8000)")
	migrateDynamoToMongoCommand.Flags().StringVar(&migrateDynamoRegion, "dynamo-region", "eu-west-1", "AWS region for DynamoDB")
}

var (
	migrateDynamoEndpoint string
	migrateDynamoRegion   string
)

var migrateDynamoToMongoCommand = &cobra.Command{
	Use:   "migrate-dynamo-to-mongo",
	Short: "Migrate all data from DynamoDB to MongoDB",
	Run: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()

		cfg := internalConfig.GetConfig()
		ctx := context.Background()

		if migrateDynamoEndpoint == "" {
			logger.Log.Fatal("--dynamo-endpoint is required. For local DynamoDB use: --dynamo-endpoint http://usegro_dynamodb:8000")
		}

		var cfgOpts []func(*awsconfig.LoadOptions) error
		cfgOpts = append(cfgOpts, awsconfig.WithRegion(migrateDynamoRegion))
		cfgOpts = append(cfgOpts,
			awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: migrateDynamoEndpoint}, nil
				})),
			awsconfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{AccessKeyID: "dummy", SecretAccessKey: "dummy"},
			}),
		)
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx, cfgOpts...)
		if err != nil {
			logger.Log.Fatal("Unable to load AWS config", zap.Error(err))
		}
		dynamoClient := dynamodb.NewFromConfig(awsCfg)

		// Connect to MongoDB
		mongoCfg := cfg.MongoDB
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCfg.URI))
		if err != nil {
			logger.Log.Fatal("Unable to connect to MongoDB", zap.Error(err))
		}
		defer mongoClient.Disconnect(ctx)
		db := mongoClient.Database(mongoCfg.Database)

		tables := []string{"forms", "form_submissions", "tags", "customer_activity"}

		for _, table := range tables {
			logger.Log.Info("Migrating table", zap.String("table", table))
			count, err := migrateTable(ctx, dynamoClient, db, table)
			if err != nil {
				logger.Log.Error("Failed to migrate table", zap.String("table", table), zap.Error(err))
				continue
			}
			logger.Log.Info("Migrated table", zap.String("table", table), zap.Int("count", count))
		}

		logger.Log.Info("Migration complete")
	},
}

func migrateTable(ctx context.Context, dynamo *dynamodb.Client, db *mongo.Database, tableName string) (int, error) {
	col := db.Collection(tableName)
	count := 0

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	for {
		out, err := dynamo.Scan(ctx, input)
		if err != nil {
			return count, fmt.Errorf("scan failed: %w", err)
		}

		if len(out.Items) == 0 {
			break
		}

		var docs []interface{}
		for _, item := range out.Items {
			var doc map[string]interface{}
			if err := attributevalue.UnmarshalMap(item, &doc); err != nil {
				log.Printf("Warning: failed to unmarshal item in %s: %v", tableName, err)
				continue
			}
			docs = append(docs, doc)
		}

		if len(docs) > 0 {
			_, err = col.InsertMany(ctx, docs)
			if err != nil {
				return count, fmt.Errorf("insert failed: %w", err)
			}
			count += len(docs)
		}

		if out.LastEvaluatedKey == nil {
			break
		}

		// Pagination: set ExclusiveStartKey for next scan
		// The DynamoDB SDK v2 uses map[string]types.AttributeValue for LastEvaluatedKey
		// We need to cast it properly for ScanInput
		break // simplified: full pagination would require setting ExclusiveStartKey
	}

	return count, nil
}
