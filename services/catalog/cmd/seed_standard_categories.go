package cmd

import (
	_ "embed"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/usegro/services/catalog/database"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/logger"
	"go.uber.org/zap"
)

//go:embed data/standard_categories.json
var standardCategoriesJSON []byte

//go:embed data/taxonomy_attributes.json
var taxonomyAttributesJSON []byte

// categoryNode mirrors the raw JSON structure (string IDs).
type categoryNode struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Children []categoryNode `json:"children"`
}

type taxonomyFile struct {
	Attributes []taxonomyAttr `json:"attributes"`
	Categories []taxonomyCat  `json:"categories"`
}

type taxonomyAttr struct {
	ID     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Handle string          `json:"handle"`
	Values []taxonomyValue `json:"values"`
}

type taxonomyValue struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Handle string    `json:"handle"`
}

type taxonomyCat struct {
	ID         uuid.UUID      `json:"id"`
	FullName   string         `json:"full_name"`
	Level      int            `json:"level"`
	Attributes []taxonomyAttr `json:"attributes"`
}

// categoryIDNamespace is the UUID v5 namespace used to derive deterministic
// UUIDs for categories from the short string IDs in standard_categories.json.
// Attribute and value UUIDs are read directly from taxonomy_attributes.json.
var categoryIDNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // uuid.NameSpaceURL

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: "seed", Title: "Seed:"})
	rootCmd.AddCommand(seedStandardCategoriesCommand)
}

var seedStandardCategoriesCommand = &cobra.Command{
	Use:     "seed:standard-categories",
	Short:   "Seed standard_categories and standard_attributes tables from the product taxonomy JSON",
	GroupID: "seed",
	Run: func(cmd *cobra.Command, _ []string) {
		setUpConfig()
		setUpLogger()
		database.SetUpPostgres()

		db := database.PostgressInstance1

		// ── 0. Parse taxonomy file ────────────────────────────────────────────
		var taxonomy taxonomyFile
		if err := json.Unmarshal(taxonomyAttributesJSON, &taxonomy); err != nil {
			logger.Log.Fatal("failed to parse taxonomy_attributes.json", zap.Error(err))
		}

		// Build lookup: category UUID → metadata (full_name, level, attributes)
		catMetaByID := make(map[uuid.UUID]taxonomyCat, len(taxonomy.Categories))
		for _, c := range taxonomy.Categories {
			catMetaByID[c.ID] = c
		}

		// ── 1. Seed categories ────────────────────────────────────────────────
		var roots []categoryNode
		if err := json.Unmarshal(standardCategoriesJSON, &roots); err != nil {
			logger.Log.Fatal("failed to parse standard_categories.json", zap.Error(err))
		}

		var categoryRecords []models.StandardCategory
		flattenNodes(roots, nil, catMetaByID, &categoryRecords)

		logger.Log.Info("seeding standard categories", zap.Int("count", len(categoryRecords)))
		for i := 0; i < len(categoryRecords); i += 500 {
			end := i + 500
			if end > len(categoryRecords) {
				end = len(categoryRecords)
			}
			batch := categoryRecords[i:end]
			if err := db.Save(&batch).Error; err != nil {
				logger.Log.Fatal("failed to upsert categories", zap.Int("offset", i), zap.Error(err))
			}
		}
		logger.Log.Info("✅ standard categories seeded", zap.Int("total", len(categoryRecords)))

		// ── 2. Seed attributes ────────────────────────────────────────────────
		attrRecords := make([]models.StandardAttribute, 0, len(taxonomy.Attributes))
		for _, a := range taxonomy.Attributes {
			values := make(models.AttributeValues, len(a.Values))
			for i, v := range a.Values {
				values[i] = models.AttributeValue{ID: v.ID.String(), Name: v.Name, Handle: v.Handle}
			}
			attrRecords = append(attrRecords, models.StandardAttribute{
				ID:     a.ID,
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
			batch := attrRecords[i:end]
			if err := db.Save(&batch).Error; err != nil {
				logger.Log.Fatal("failed to upsert attributes", zap.Int("offset", i), zap.Error(err))
			}
		}
		logger.Log.Info("✅ standard attributes seeded", zap.Int("total", len(attrRecords)))

		// ── 3. Link categories → attributes ──────────────────────────────────
		type joinRow struct {
			StandardCategoryID  uuid.UUID
			StandardAttributeID uuid.UUID
		}

		var joins []joinRow
		for _, cat := range taxonomy.Categories {
			for _, a := range cat.Attributes {
				joins = append(joins, joinRow{
					StandardCategoryID:  cat.ID,
					StandardAttributeID: a.ID,
				})
			}
		}

		logger.Log.Info("seeding category-attribute links", zap.Int("count", len(joins)))
		for i := 0; i < len(joins); i += 500 {
			end := i + 500
			if end > len(joins) {
				end = len(joins)
			}
			if err := db.Table("standard_category_attributes").Save(joins[i:end]).Error; err != nil {
				logger.Log.Fatal("failed to upsert category-attribute links", zap.Int("offset", i), zap.Error(err))
			}
		}
		logger.Log.Info("✅ category-attribute links seeded", zap.Int("total", len(joins)))
	},
}

// flattenNodes recursively converts the category tree into a flat slice of
// StandardCategory records. UUIDs are derived from the short string IDs in
// standard_categories.json using the same namespace as taxonomy_attributes.json,
// so they match the UUIDs stored in that file.
func flattenNodes(nodes []categoryNode, parentID *uuid.UUID, metaByID map[uuid.UUID]taxonomyCat, out *[]models.StandardCategory) {
	for _, n := range nodes {
		id := uuid.NewSHA1(categoryIDNamespace, []byte(n.ID))

		record := models.StandardCategory{
			ID:       id,
			ParentID: parentID,
			Name:     n.Name,
			IsLeaf:   len(n.Children) == 0,
		}

		if meta, ok := metaByID[id]; ok {
			record.FullName = meta.FullName
			record.Level = meta.Level
		}

		*out = append(*out, record)
		if len(n.Children) > 0 {
			flattenNodes(n.Children, &id, metaByID, out)
		}
	}
}
