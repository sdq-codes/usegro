package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	internalConfig "github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(reformatMongoCommand)
}

var reformatMongoCommand = &cobra.Command{
	Use:   "reformat-mongo",
	Short: "Reformat MongoDB collections from DynamoDB single-table layout to MongoDB-native layout",
	Run: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()

		cfg := internalConfig.GetConfig()
		ctx := context.Background()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
		if err != nil {
			logger.Log.Fatal("Unable to connect to MongoDB", zap.Error(err))
		}
		defer client.Disconnect(ctx)
		db := client.Database(cfg.MongoDB.Database)

		logger.Log.Info("Starting MongoDB reformat")

		// Drop stale PK/SK indexes from the DynamoDB era — these block upserts
		// because null PK+SK would violate the unique constraint.
		for _, colName := range []string{"forms", "form_submissions", "tags", "customer_activity"} {
			col := db.Collection(colName)
			if _, err := col.Indexes().DropOne(ctx, "PK_1_SK_1"); err != nil {
				// Index may not exist — that's fine
				logger.Log.Debug("Drop PK_1_SK_1 index (may not exist)", zap.String("collection", colName), zap.Error(err))
			} else {
				logger.Log.Info("Dropped PK_1_SK_1 index", zap.String("collection", colName))
			}
		}

		if n, err := reformatForms(ctx, db); err != nil {
			logger.Log.Error("Failed to reformat forms", zap.Error(err))
		} else {
			logger.Log.Info("Reformatted forms collection", zap.Int("processed", n))
		}

		if n, err := reformatSubmissions(ctx, db); err != nil {
			logger.Log.Error("Failed to reformat submissions", zap.Error(err))
		} else {
			logger.Log.Info("Reformatted form_submissions collection", zap.Int("processed", n))
		}

		if n, err := reformatTags(ctx, db); err != nil {
			logger.Log.Error("Failed to reformat tags", zap.Error(err))
		} else {
			logger.Log.Info("Reformatted tags collection", zap.Int("processed", n))
		}

		if n, err := reformatActivity(ctx, db); err != nil {
			logger.Log.Error("Failed to reformat customer_activity", zap.Error(err))
		} else {
			logger.Log.Info("Reformatted customer_activity collection", zap.Int("processed", n))
		}

		logger.Log.Info("Reformat complete")
	},
}

// reformatForms handles the single-table DynamoDB dump in the `forms` collection.
// Items are classified by their SK value and split into forms, form_versions, form_fields.
func reformatForms(ctx context.Context, db *mongo.Database) (int, error) {
	col := db.Collection("forms")
	versionsCol := db.Collection("form_versions")
	fieldsCol := db.Collection("form_fields")

	cur, err := col.Find(ctx, bson.M{"PK": bson.M{"$exists": true}})
	if err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}
	defer cur.Close(ctx)

	count := 0
	for cur.Next(ctx) {
		var raw bson.M
		if err := cur.Decode(&raw); err != nil {
			logger.Log.Warn("Failed to decode forms doc", zap.Error(err))
			continue
		}

		pk, _ := raw["PK"].(string)
		sk, _ := raw["SK"].(string)
		oldID := raw["_id"]

		if pk == "" || sk == "" {
			continue
		}

		formID := strings.TrimPrefix(pk, "FORM#")
		delete(raw, "PK")
		delete(raw, "SK")
		delete(raw, "_id")

		switch {
		case sk == "METADATA":
			// → forms collection
			raw["_id"] = formID
			if _, err := col.ReplaceOne(ctx,
				bson.M{"_id": formID},
				raw,
				options.Replace().SetUpsert(true),
			); err != nil {
				logger.Log.Warn("Failed to upsert form", zap.String("formID", formID), zap.Error(err))
				continue
			}
			// Delete the original PK/SK document (different _id)
			if oldID != nil {
				col.DeleteOne(ctx, bson.M{"_id": oldID})
			}

		case strings.HasPrefix(sk, "VERSION#") && !strings.Contains(sk, "FIELD#"):
			// → form_versions collection
			versionID := strings.TrimPrefix(sk, "VERSION#")
			raw["_id"] = versionID
			raw["formID"] = formID
			if _, err := versionsCol.ReplaceOne(ctx,
				bson.M{"_id": versionID},
				raw,
				options.Replace().SetUpsert(true),
			); err != nil {
				logger.Log.Warn("Failed to upsert form_version", zap.String("versionID", versionID), zap.Error(err))
				continue
			}
			if oldID != nil {
				col.DeleteOne(ctx, bson.M{"_id": oldID})
			}

		case strings.Contains(sk, "FIELD#"):
			// → form_fields collection
			// SK format: "VERSION#<versionID>FIELD#<fieldID>"
			parts := strings.SplitN(sk, "FIELD#", 2)
			if len(parts) != 2 {
				logger.Log.Warn("Unexpected field SK format", zap.String("sk", sk))
				continue
			}
			versionID := strings.TrimPrefix(parts[0], "VERSION#")
			fieldID := parts[1]

			raw["_id"] = fieldID
			// Normalise formVersionID — old data stored the full SK "VERSION#<id>", strip to bare UUID
			if fvID, ok := raw["formVersionID"].(string); ok && strings.HasPrefix(fvID, "VERSION#") {
				raw["formVersionID"] = strings.TrimPrefix(fvID, "VERSION#")
			} else {
				raw["formVersionID"] = versionID
			}

			if _, err := fieldsCol.ReplaceOne(ctx,
				bson.M{"_id": fieldID},
				raw,
				options.Replace().SetUpsert(true),
			); err != nil {
				logger.Log.Warn("Failed to upsert form_field", zap.String("fieldID", fieldID), zap.Error(err))
				continue
			}
			if oldID != nil {
				col.DeleteOne(ctx, bson.M{"_id": oldID})
			}

		default:
			logger.Log.Warn("Unknown SK pattern in forms collection", zap.String("sk", sk))
			continue
		}

		count++
	}
	return count, cur.Err()
}

// reformatSubmissions rewrites form_submissions documents from PK/SK format to _id=submissionID.
func reformatSubmissions(ctx context.Context, db *mongo.Database) (int, error) {
	col := db.Collection("form_submissions")

	cur, err := col.Find(ctx, bson.M{"PK": bson.M{"$exists": true}})
	if err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}
	defer cur.Close(ctx)

	count := 0
	for cur.Next(ctx) {
		var raw bson.M
		if err := cur.Decode(&raw); err != nil {
			logger.Log.Warn("Failed to decode submission doc", zap.Error(err))
			continue
		}

		oldID := raw["_id"]
		submissionID, _ := raw["submissionID"].(string)
		if submissionID == "" {
			// Fall back: extract from SK "SUBMISSION#<id>"
			if sk, ok := raw["SK"].(string); ok {
				submissionID = strings.TrimPrefix(sk, "SUBMISSION#")
			}
		}
		if submissionID == "" {
			logger.Log.Warn("Cannot determine submissionID, skipping", zap.Any("_id", oldID))
			continue
		}

		delete(raw, "PK")
		delete(raw, "SK")
		delete(raw, "submissionID") // was a separate field; now it lives as _id
		delete(raw, "_id")
		raw["_id"] = submissionID

		// Ensure status field exists
		if _, ok := raw["status"]; !ok {
			raw["status"] = "active"
		}

		if _, err := col.ReplaceOne(ctx,
			bson.M{"_id": submissionID},
			raw,
			options.Replace().SetUpsert(true),
		); err != nil {
			logger.Log.Warn("Failed to upsert submission", zap.String("submissionID", submissionID), zap.Error(err))
			continue
		}
		if oldID != nil && oldID != submissionID {
			col.DeleteOne(ctx, bson.M{"_id": oldID})
		}
		count++
	}
	return count, cur.Err()
}

// reformatTags rewrites tags documents from PK/SK format to _id=tagID.
func reformatTags(ctx context.Context, db *mongo.Database) (int, error) {
	col := db.Collection("tags")

	cur, err := col.Find(ctx, bson.M{"PK": bson.M{"$exists": true}})
	if err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}
	defer cur.Close(ctx)

	count := 0
	for cur.Next(ctx) {
		var raw bson.M
		if err := cur.Decode(&raw); err != nil {
			logger.Log.Warn("Failed to decode tag doc", zap.Error(err))
			continue
		}

		oldID := raw["_id"]
		sk, _ := raw["SK"].(string)
		pk, _ := raw["PK"].(string)

		tagID := strings.TrimPrefix(sk, "TAG#")
		if tagID == "" || tagID == sk {
			logger.Log.Warn("Cannot determine tagID from SK", zap.String("sk", sk))
			continue
		}

		// crmID may already be a field; if not, extract from PK
		if _, ok := raw["crmID"]; !ok {
			raw["crmID"] = strings.TrimPrefix(pk, "CRM#")
		}

		delete(raw, "PK")
		delete(raw, "SK")
		delete(raw, "_id")
		raw["_id"] = tagID

		if _, err := col.ReplaceOne(ctx,
			bson.M{"_id": tagID},
			raw,
			options.Replace().SetUpsert(true),
		); err != nil {
			logger.Log.Warn("Failed to upsert tag", zap.String("tagID", tagID), zap.Error(err))
			continue
		}
		if oldID != nil && oldID != tagID {
			col.DeleteOne(ctx, bson.M{"_id": oldID})
		}
		count++
	}
	return count, cur.Err()
}

// reformatActivity assigns proper _id UUIDs to customer_activity docs that lack one.
func reformatActivity(ctx context.Context, db *mongo.Database) (int, error) {
	col := db.Collection("customer_activity")

	cur, err := col.Find(ctx, bson.M{"PK": bson.M{"$exists": true}})
	if err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}
	defer cur.Close(ctx)

	count := 0
	for cur.Next(ctx) {
		var raw bson.M
		if err := cur.Decode(&raw); err != nil {
			logger.Log.Warn("Failed to decode activity doc", zap.Error(err))
			continue
		}

		oldID := raw["_id"]
		pk, _ := raw["PK"].(string)

		if _, ok := raw["crmID"]; !ok {
			raw["crmID"] = strings.TrimPrefix(pk, "CRM#")
		}

		newID := uuid.New().String()
		delete(raw, "PK")
		delete(raw, "SK")
		delete(raw, "_id")
		raw["_id"] = newID

		if _, err := col.InsertOne(ctx, raw); err != nil {
			logger.Log.Warn("Failed to insert activity", zap.Error(err))
			continue
		}
		if oldID != nil {
			col.DeleteOne(ctx, bson.M{"_id": oldID})
		}
		count++
	}
	return count, cur.Err()
}
