package migrations

import (
	baseModels "github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"go.uber.org/zap"
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

	if err := db.AutoMigrate(
		&baseModels.User{},
	); err != nil {
		logger.Log.Error("Automigration failed", zap.Error(err))
		return err
	}

	logger.Log.Info("✅ AutoMigrate completed.")
	return nil
}
