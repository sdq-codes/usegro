package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/usegro/services/catalog/database"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

// attrIDNamespace is the UUID v5 namespace for deriving deterministic attribute UUIDs
// from Shopify taxonomy attribute handles.
var attrIDNamespace = uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890")

type shopifyMetafields struct {
	Attributes []shopifyAttribute `json:"attributes"`
	Categories []shopifyCategory  `json:"categories"`
}

type shopifyAttribute struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Handle string         `json:"handle"`
	Values []shopifyValue `json:"values"`
}

type shopifyValue struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Handle string `json:"handle"`
}

type shopifyCategory struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	FullName   string             `json:"full_name"`
	Level      int                `json:"level"`
	Attributes []shopifyAttribute `json:"attributes"`
}

func init() {
	seedShopifyCommand.Flags().String("file", "", "Path to shopify_metafields.json (required)")
	_ = seedShopifyCommand.MarkFlagRequired("file")
	rootCmd.AddCommand(seedShopifyCommand)
}

var seedShopifyCommand = &cobra.Command{
	Use:     "seed:shopify-categories",
	Short:   "Seed standard_categories and standard_attributes from shopify_metafields.json",
	GroupID: "seed",
	Run: func(cmd *cobra.Command, _ []string) {
		filePath, _ := cmd.Flags().GetString("file")

		setUpConfig()
		setUpLogger()
		database.SetUpPostgres()
		db := database.PostgressInstance1

		// ── 0. Parse file ─────────────────────────────────────────────────────
		logger.Log.Info("reading shopify_metafields.json", zap.String("path", filePath))
		raw, err := os.ReadFile(filePath)
		if err != nil {
			logger.Log.Fatal("failed to read file", zap.Error(err))
		}

		var meta shopifyMetafields
		if err := json.Unmarshal(raw, &meta); err != nil {
			logger.Log.Fatal("failed to parse shopify_metafields.json", zap.Error(err))
		}
		logger.Log.Info("parsed file",
			zap.Int("categories", len(meta.Categories)),
			zap.Int("attributes", len(meta.Attributes)))

		// ── 1. Compute is_leaf (not present in GitHub-sourced data) ───────────
		// A category is a leaf if no other category has it as its parent.
		parentSet := make(map[string]struct{}, len(meta.Categories))
		for _, c := range meta.Categories {
			shortID := shopifyCatShortID(c.ID)
			if p := shopifyParentShortID(shortID); p != "" {
				parentSet[p] = struct{}{}
			}
		}

		// ── 2. Seed categories ─────────────────────────────────────────────────
		categoryRecords := make([]models.StandardCategory, 0, len(meta.Categories))
		for _, c := range meta.Categories {
			shortID := shopifyCatShortID(c.ID)
			id := uuid.NewSHA1(categoryIDNamespace, []byte(shortID))

			var parentID *uuid.UUID
			if p := shopifyParentShortID(shortID); p != "" {
				pid := uuid.NewSHA1(categoryIDNamespace, []byte(p))
				parentID = &pid
			}

			_, hasChildren := parentSet[shortID]

			categoryRecords = append(categoryRecords, models.StandardCategory{
				ID:       id,
				ParentID: parentID,
				Name:     c.Name,
				FullName: c.FullName,
				Level:    c.Level,
				IsLeaf:   !hasChildren,
			})
		}

		logger.Log.Info("seeding standard categories", zap.Int("count", len(categoryRecords)))
		for i := 0; i < len(categoryRecords); i += 500 {
			end := i + 500
			if end > len(categoryRecords) {
				end = len(categoryRecords)
			}
			if err := db.Save(categoryRecords[i:end]).Error; err != nil {
				logger.Log.Fatal("failed to upsert categories", zap.Int("offset", i), zap.Error(err))
			}
		}
		logger.Log.Info("✅ standard categories seeded", zap.Int("total", len(categoryRecords)))

		// ── 3. Seed attributes ─────────────────────────────────────────────────
		attrRecords := make([]models.StandardAttribute, 0, len(meta.Attributes))
		for _, a := range meta.Attributes {
			if a.Handle == "" {
				continue
			}
			id := uuid.NewSHA1(attrIDNamespace, []byte(a.Handle))
			values := make(models.AttributeValues, len(a.Values))
			for i, v := range a.Values {
				values[i] = models.AttributeValue{
					ID:     v.ID,
					Name:   v.Name,
					Handle: v.Handle,
				}
			}
			attrRecords = append(attrRecords, models.StandardAttribute{
				ID:     id,
				Name:   a.Name,
				Handle: a.Handle,
				Values: values,
			})
		}

		logger.Log.Info("seeding standard attributes", zap.Int("count", len(attrRecords)))
		for i := 0; i < len(attrRecords); i += 500 {
			end := i + 500
			if end > len(attrRecords) {
				end = len(attrRecords)
			}
			if err := db.Save(attrRecords[i:end]).Error; err != nil {
				logger.Log.Fatal("failed to upsert attributes", zap.Int("offset", i), zap.Error(err))
			}
		}
		logger.Log.Info("✅ standard attributes seeded", zap.Int("total", len(attrRecords)))

		// ── 4. Link categories → attributes ───────────────────────────────────
		type joinRow struct {
			StandardCategoryID  uuid.UUID `gorm:"column:standard_category_id"`
			StandardAttributeID uuid.UUID `gorm:"column:standard_attribute_id"`
		}

		var joins []joinRow
		for _, cat := range meta.Categories {
			if len(cat.Attributes) == 0 {
				continue
			}
			catUUID := uuid.NewSHA1(categoryIDNamespace, []byte(shopifyCatShortID(cat.ID)))
			for _, a := range cat.Attributes {
				if a.Handle == "" {
					continue
				}
				attrUUID := uuid.NewSHA1(attrIDNamespace, []byte(a.Handle))
				joins = append(joins, joinRow{
					StandardCategoryID:  catUUID,
					StandardAttributeID: attrUUID,
				})
			}
		}

		logger.Log.Info("seeding category-attribute links", zap.Int("count", len(joins)))
		for i := 0; i < len(joins); i += 500 {
			end := i + 500
			if end > len(joins) {
				end = len(joins)
			}
			if err := db.Table("standard_category_attributes").
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(joins[i:end]).Error; err != nil {
				logger.Log.Fatal("failed to insert category-attribute links", zap.Int("offset", i), zap.Error(err))
			}
		}
		logger.Log.Info("✅ category-attribute links seeded", zap.Int("total", len(joins)))
	},
}

// shopifyCatShortID extracts the short ID from a Shopify taxonomy category GID.
// e.g. "gid://shopify/TaxonomyCategory/aa-1-2" → "aa-1-2"
func shopifyCatShortID(gid string) string {
	i := strings.LastIndex(gid, "/")
	if i < 0 {
		return gid
	}
	return gid[i+1:]
}

// shopifyParentShortID returns the parent short ID by stripping the last "-N" segment.
// e.g. "aa-1-2" → "aa-1", "aa-1" → "aa", "aa" → ""
func shopifyParentShortID(shortID string) string {
	i := strings.LastIndex(shortID, "-")
	if i < 0 {
		return ""
	}
	return shortID[:i]
}
