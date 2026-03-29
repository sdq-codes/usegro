package database

import (
	"context"
	"fmt"
	"time"

	"github.com/usegro/services/catalog/config"
	catalogModels "github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgressInstance1 *gorm.DB

func SetUpPostgres() {
	cfg := config.GetConfig().Postgres

	if cfg.Host == "" {
		logger.Log.Fatal("Postgres host is empty — skipping setup")
	}

	if cfg.Schema == "" {
		logger.Log.Fatal("Postgres schema is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.Port,
		sslMode,
	)), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("gorm.Open() failed", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Fatal("db.DB() failed", zap.Error(err))
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Log.Fatal("Postgres ping failed", zap.Error(err))
	}

	logger.Log.Info("✅ GORM connection established")

	if err := db.AutoMigrate(
		&catalogModels.StandardAttribute{},
		&catalogModels.StandardCategory{},
		&catalogModels.Category{},
		&catalogModels.CatalogTag{},
		&catalogModels.CatalogItem{},
		&catalogModels.CatalogAdditionalField{},
		&catalogModels.ProductDetail{},
		&catalogModels.ProductOption{},
		&catalogModels.ProductOptionValue{},
		&catalogModels.ProductVariant{},
		&catalogModels.ProductVariantGroup{},
		&catalogModels.ServiceDetail{},
		&catalogModels.ServiceVariation{},
		&catalogModels.ServiceLocation{},
		&catalogModels.Plan{},
		&catalogModels.CustomerPlan{},
		&catalogModels.Booking{},
		&catalogModels.CatalogItemMedia{},
	); err != nil {
		logger.Log.Fatal("AutoMigrate failed", zap.Error(err))
	}
	logger.Log.Info("✅ AutoMigration complete")

	PostgressInstance1 = db
}
