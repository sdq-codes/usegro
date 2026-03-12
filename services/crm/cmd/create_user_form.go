package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	dynamodb2 "github.com/usegro/services/crm/internal/helper/dynamodb"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	internalConfig "github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/logger"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(seedCommand)

	seedCommand.Flags().StringP("file", "f", "./database/seed/dynamo/createCustomerForm.json", "(required) Path to JSON seed file")
	seedCommand.Flags().StringP("table", "t", "", "DynamoDB table name (uses env var if not provided)")
	seedCommand.Flags().StringP("crm-id", "c", "", "CRM ID (generates UUID if not provided)")
}

var seedCommand = &cobra.Command{
	Use:     "seed",
	Short:   "Seed data into DynamoDB from JSON file",
	GroupID: "make",
	Run: func(cmd *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()

		filePath, _ := cmd.Flags().GetString("file")
		crmID := "global"

		dynamoCfg := internalConfig.GetConfig().DynamodbForms
		tableName := dynamoCfg.DynamoFormTableName
		if tableName == "" {
			logger.Log.Fatal("Table name not configured (dynamodbForms.dynamoFormTableName)")
		}

		ctx := context.Background()

		var cfgOpts []func(*config.LoadOptions) error
		cfgOpts = append(cfgOpts, config.WithRegion(dynamoCfg.AwsRegion))

		if dynamoCfg.DynamoEndpoint != "" {
			// Local DynamoDB
			cfgOpts = append(cfgOpts,
				config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{URL: dynamoCfg.DynamoEndpoint}, nil
					})),
				config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
					Value: aws.Credentials{
						AccessKeyID: "dummy", SecretAccessKey: "dummy",
					},
				}),
			)
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			logger.Log.Fatal("Unable to load AWS SDK config", zap.Error(err))
		}

		client := dynamodb.NewFromConfig(cfg)

		// Read and parse JSON file
		logger.Log.Info("Reading seed file", zap.String("file", filePath))
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			logger.Log.Fatal("Error reading JSON file", zap.Error(err))
		}

		var jsonFields []JSONField
		if err := json.Unmarshal(jsonData, &jsonFields); err != nil {
			logger.Log.Fatal("Error parsing JSON", zap.Error(err))
		}

		if len(jsonFields) == 0 {
			logger.Log.Fatal("No data found in JSON file")
		}

		// Generate IDs
		formID := uuid.New().String()
		now := time.Now()

		logger.Log.Info("Starting seed process",
			zap.String("crmID", crmID),
			zap.String("formID", formID),
			zap.String("table", tableName),
		)

		// Create Form
		form := Form{
			PK:        fmt.Sprintf("FORM#%s", formID),
			SK:        "METADATA",
			CrmID:     crmID,
			Type:      jsonFields[0].Type,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := putItem(ctx, client, tableName, form); err != nil {
			logger.Log.Fatal("Error creating form", zap.Error(err))
		}
		logger.Log.Info("Created form", zap.String("formID", formID))

		versionID := uuid.New().String()
		// Create FormVersion
		formVersion := FormVersion{
			PK:                fmt.Sprintf("FORM#%s", formID),
			SK:                fmt.Sprintf("VERSION#%s", versionID),
			FormID:            formID,
			Title:             jsonFields[0].Title,
			Description:       jsonFields[0].Description,
			FormVersionStatus: "published",
			PublishedAt:       now.String(),
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		if err := putItem(ctx, client, tableName, formVersion); err != nil {
			logger.Log.Fatal("Error creating form version", zap.Error(err))
		}
		logger.Log.Info("Created form version", zap.String("formVersionID", versionID))

		// Create FormVersionFields and FieldLogic
		fieldIDMap := make(map[string]string) // slug -> fieldID
		fieldSKMap := make(map[string]string) // slug -> SK
		fieldCount := 0

		for _, jsonField := range jsonFields[1:] {
			fieldID := uuid.New().String()
			fieldSK := fmt.Sprintf("VERSION#%sFIELD#%s", versionID, fieldID)

			fieldIDMap[jsonField.Slug] = fieldID
			fieldSKMap[jsonField.Slug] = fieldSK

			field := FormVersionField{
				PK:            fmt.Sprintf("FORM#%s", dynamodb2.ExtractSK(formVersion.PK)),
				SK:            fmt.Sprintf("VERSION#%sFIELD#%s", dynamodb2.ExtractSK(formVersion.SK), fieldID),
				FormVersionID: fmt.Sprintf("VERSION#%s", versionID),
				FieldTypeID:   jsonField.FieldTypeID,
				FieldTypeName: jsonField.FieldTypeName,
				Label:         jsonField.Label,
				Description:   jsonField.Description,
				Hint:          jsonField.Hint,
				Section:       jsonField.Section,
				Placeholder:   jsonField.Placeholder,
				Configs:       jsonField.Configs,
				Options:       jsonField.Options,
				Alert:         jsonField.Alert,
				Validations:   jsonField.Validations,
				Order:         jsonField.Order,
				Required:      jsonField.Required,
				Slug:          jsonField.Slug,
				Logic:         []FieldLogic{},
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			if err := putItem(ctx, client, tableName, field); err != nil {
				logger.Log.Fatal("Error creating field",
					zap.String("slug", jsonField.Slug),
					zap.Error(err),
				)
			}
			fieldCount++
			logger.Log.Info("Created field",
				zap.String("slug", jsonField.Slug),
				zap.String("fieldID", fieldID),
				zap.String("section", jsonField.Section),
			)
		}
	},
}

// Models
type Form struct {
	PK        string    `dynamodbav:"PK"`
	SK        string    `dynamodbav:"SK"`
	CrmID     string    `json:"crmID" dynamodbav:"crmID"`
	Type      string    `json:"type" dynamodbav:"type"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type FormVersion struct {
	PK                string    `dynamodbav:"PK"`
	SK                string    `dynamodbav:"SK"`
	FormID            string    `json:"formID" dynamodbav:"formID"`
	Title             string    `json:"title" dynamodbav:"title"`
	Description       string    `json:"description" dynamodbav:"description"`
	FormVersionStatus string    `json:"formVersionStatus" dynamodbav:"formVersionStatus"`
	PublishedAt       string    `json:"publishedAt" dynamodbav:"publishedAt"`
	CreatedAt         time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type Option struct {
	Label string `dynamodbav:"label" json:"label"`
	Value string `dynamodbav:"value" json:"value"`
}

type Alert struct {
	Icon    string `dynamodbav:"icon" json:"icon"`
	Type    string `dynamodbav:"type" json:"type"`
	Message string `dynamodbav:"message" json:"message"`
}

type FormVersionField struct {
	PK            string                   `dynamodbav:"PK"`
	SK            string                   `dynamodbav:"SK"`
	FormVersionID string                   `dynamodbav:"formVersionID" json:"formVersionID"`
	FieldTypeID   uint                     `dynamodbav:"fieldTypeId" json:"fieldTypeID"`
	FieldTypeName string                   `dynamodbav:"fieldTypeName" json:"fieldTypeName"`
	Label         string                   `dynamodbav:"label" json:"label"`
	Description   string                   `dynamodbav:"description" json:"description"`
	Hint          string                   `dynamodbav:"hint" json:"hint"`
	Section       string                   `dynamodbav:"section" json:"section"`
	Placeholder   string                   `dynamodbav:"placeholder" json:"placeholder"`
	Configs       []map[string]interface{} `dynamodbav:"configs" json:"configs"`
	Options       []Option                 `dynamodbav:"options" json:"options"`
	Alert         []Alert                  `dynamodbav:"alert" json:"alert"`
	Validations   []map[string]string      `dynamodbav:"validations" json:"validations"`
	Order         int                      `dynamodbav:"order" json:"order"`
	Required      bool                     `dynamodbav:"required" json:"required"`
	Slug          string                   `dynamodbav:"slug" json:"slug"`
	Logic         []FieldLogic             `dynamodbav:"logic" json:"logic"`
	CreatedAt     time.Time                `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time                `dynamodbav:"updatedAt" json:"updatedAt"`
}

type FieldLogic struct {
	PK                 string      `dynamodbav:"PK"`
	SK                 string      `dynamodbav:"SK"`
	FormVersionFieldID string      `dynamodbav:"formVersionFieldID" json:"formVersionFieldID"`
	Operator           string      `dynamodbav:"operator" json:"operator"`
	Value              interface{} `dynamodbav:"value" json:"value"`
	Action             string      `dynamodbav:"action" json:"action"`
	CreatedAt          time.Time   `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time   `dynamodbav:"updatedAt" json:"updatedAt"`
}

type JSONField struct {
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	Placeholder   string                   `json:"placeholder"`
	Type          string                   `json:"type"`
	Label         string                   `json:"label"`
	Section       string                   `json:"section"`
	Hint          string                   `json:"hint"`
	Order         int                      `json:"order"`
	FieldTypeID   uint                     `json:"fieldTypeID"`
	FieldTypeName string                   `json:"fieldTypeName"`
	Validations   []map[string]string      `json:"validations"`
	Configs       []map[string]interface{} `json:"configs"`
	Slug          string                   `json:"slug"`
	Options       []Option                 `json:"options"`
	Required      bool                     `json:"required"`
	Alert         []Alert                  `json:"alert"`
}

// Helper function
func putItem(ctx context.Context, client *dynamodb.Client, tableName string, item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}
