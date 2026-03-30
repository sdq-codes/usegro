package cmd

import (
	"context"

	"github.com/spf13/cobra"
	internalConfig "github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(createTablesCommand)
}

var createTablesCommand = &cobra.Command{
	Use:     "create-tables",
	Short:   "Create all required MongoDB indexes",
	GroupID: "make",
	Run: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()

		mongoCfg := internalConfig.GetConfig().MongoDB
		ctx := context.Background()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCfg.URI))
		if err != nil {
			logger.Log.Fatal("Unable to connect to MongoDB", zap.Error(err))
		}
		defer client.Disconnect(ctx)
		db := client.Database(mongoCfg.Database)

		indexDefs := map[string][]mongo.IndexModel{
			"forms": {
				{Keys: bson.D{{Key: "crmID", Value: 1}, {Key: "type", Value: 1}}},
			},
			"form_versions": {
				{Keys: bson.D{{Key: "formID", Value: 1}}},
				{Keys: bson.D{{Key: "formID", Value: 1}, {Key: "formVersionStatus", Value: 1}}},
			},
			"form_fields": {
				{Keys: bson.D{{Key: "formVersionID", Value: 1}, {Key: "order", Value: 1}}},
			},
			"form_submissions": {
				{Keys: bson.D{{Key: "formID", Value: 1}}},
				{Keys: bson.D{{Key: "crmID", Value: 1}, {Key: "type", Value: 1}}},
				{Keys: bson.D{{Key: "formID", Value: 1}, {Key: "answers.email", Value: 1}}},
				{Keys: bson.D{{Key: "formID", Value: 1}, {Key: "answers.phone_number", Value: 1}}},
			},
			"tags": {
				{Keys: bson.D{{Key: "crmID", Value: 1}}},
			},
			"customer_activity": {
				{Keys: bson.D{{Key: "customerID", Value: 1}, {Key: "createdAt", Value: -1}}},
			},
		}

		for colName, indexes := range indexDefs {
			_, err := db.Collection(colName).Indexes().CreateMany(ctx, indexes)
			if err != nil {
				logger.Log.Error("Failed to create indexes", zap.String("collection", colName), zap.Error(err))
				continue
			}
			logger.Log.Info("Created indexes", zap.String("collection", colName))
		}

		logger.Log.Info("All indexes created successfully")
	},
}
