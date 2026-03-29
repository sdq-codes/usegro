package migrations

import (
	"github.com/usegro/services/catalog/internal/logger"
	"gorm.io/gorm"
)

type Migration struct {
	Name string
	Up   func() error
	Down func() error
}

var Migrations []*Migration

func AutoMigrateDB(db *gorm.DB) error {
	logger.Log.Info("🔄 Running AutoMigrate...")

	// Catalog-level auto-migration is handled in database.SetUpPostgres().
	// This function is reserved for manual migration steps.

	logger.Log.Info("✅ AutoMigrate completed.")
	return nil
}
