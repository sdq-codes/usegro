package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	internalConfig "github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(seedCommand)

	seedCommand.Flags().StringP("file", "f", "./database/seed/dynamo/createCustomerForm.json", "(required) Path to JSON seed file")
	seedCommand.Flags().StringP("table", "t", "", "MongoDB collection name (uses 'forms' if not provided)")
	seedCommand.Flags().StringP("crm-id", "c", "", "CRM ID (generates UUID if not provided)")
}

var seedCommand = &cobra.Command{
	Use:     "seed",
	Short:   "Seed data into MongoDB from JSON file",
	GroupID: "make",
	Run: func(cmd *cobra.Command, _ []string) {
		// Setup all the required dependencies
		setUpConfig()
		setUpLogger()

		filePath, _ := cmd.Flags().GetString("file")
		crmID := "global"

		mongoCfg := internalConfig.GetConfig().MongoDB
		ctx := context.Background()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCfg.URI))
		if err != nil {
			logger.Log.Fatal("Unable to connect to MongoDB", zap.Error(err))
		}
		defer client.Disconnect(ctx)
		db := client.Database(mongoCfg.Database)

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
		versionID := uuid.New().String()
		now := time.Now()

		logger.Log.Info("Starting seed process",
			zap.String("crmID", crmID),
			zap.String("formID", formID),
		)

		// Create Form
		form := SeedForm{
			ID:        formID,
			CrmID:     crmID,
			Type:      jsonFields[0].Type,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if _, err := db.Collection("forms").InsertOne(ctx, form); err != nil {
			logger.Log.Fatal("Error creating form", zap.Error(err))
		}
		logger.Log.Info("Created form", zap.String("formID", formID))

		// Create FormVersion
		formVersion := SeedFormVersion{
			ID:                versionID,
			FormID:            formID,
			Title:             jsonFields[0].Title,
			Description:       jsonFields[0].Description,
			FormVersionStatus: "published",
			PublishedAt:       now.String(),
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		if _, err := db.Collection("form_versions").InsertOne(ctx, formVersion); err != nil {
			logger.Log.Fatal("Error creating form version", zap.Error(err))
		}
		logger.Log.Info("Created form version", zap.String("formVersionID", versionID))

		// Create FormVersionFields
		fieldCount := 0

		for _, jsonField := range jsonFields[1:] {
			fieldID := uuid.New().String()

			field := SeedFormVersionField{
				ID:            fieldID,
				FormVersionID: versionID,
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
				Logic:         []SeedFieldLogic{},
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			if _, err := db.Collection("form_fields").InsertOne(ctx, field); err != nil {
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

		logger.Log.Info("Seed complete", zap.Int("fields", fieldCount))
	},
}

// Models
type SeedForm struct {
	ID        string    `bson:"_id"`
	CrmID     string    `json:"crmID" bson:"crmID"`
	Type      string    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SeedFormVersion struct {
	ID                string    `bson:"_id"`
	FormID            string    `json:"formID" bson:"formID"`
	Title             string    `json:"title" bson:"title"`
	Description       string    `json:"description" bson:"description"`
	FormVersionStatus string    `json:"formVersionStatus" bson:"formVersionStatus"`
	PublishedAt       string    `json:"publishedAt" bson:"publishedAt"`
	CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SeedOption struct {
	Label string `bson:"label" json:"label"`
	Value string `bson:"value" json:"value"`
}

type SeedAlert struct {
	Icon    string `bson:"icon" json:"icon"`
	Type    string `bson:"type" json:"type"`
	Message string `bson:"message" json:"message"`
}

type SeedFormVersionField struct {
	ID            string                   `bson:"_id"`
	FormVersionID string                   `bson:"formVersionID" json:"formVersionID"`
	FieldTypeID   uint                     `bson:"fieldTypeId" json:"fieldTypeID"`
	FieldTypeName string                   `bson:"fieldTypeName" json:"fieldTypeName"`
	Label         string                   `bson:"label" json:"label"`
	Description   string                   `bson:"description" json:"description"`
	Hint          string                   `bson:"hint" json:"hint"`
	Section       string                   `bson:"section" json:"section"`
	Placeholder   string                   `bson:"placeholder" json:"placeholder"`
	Configs       []map[string]interface{} `bson:"configs" json:"configs"`
	Options       []SeedOption             `bson:"options" json:"options"`
	Alert         []SeedAlert              `bson:"alert" json:"alert"`
	Validations   []map[string]string      `bson:"validations" json:"validations"`
	Order         int                      `bson:"order" json:"order"`
	Required      bool                     `bson:"required" json:"required"`
	Slug          string                   `bson:"slug" json:"slug"`
	Logic         []SeedFieldLogic         `bson:"logic" json:"logic"`
	CreatedAt     time.Time                `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time                `bson:"updatedAt" json:"updatedAt"`
}

type SeedFieldLogic struct {
	FormVersionFieldID string      `bson:"formVersionFieldID" json:"formVersionFieldID"`
	Operator           string      `bson:"operator" json:"operator"`
	Value              interface{} `bson:"value" json:"value"`
	Action             string      `bson:"action" json:"action"`
	CreatedAt          time.Time   `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time   `bson:"updatedAt" json:"updatedAt"`
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
	Options       []SeedOption             `json:"options"`
	Required      bool                     `json:"required"`
	Alert         []SeedAlert              `json:"alert"`
}
