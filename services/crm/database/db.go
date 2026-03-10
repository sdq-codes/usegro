package database

import (
	"context"
	"fmt"
	"time"

	"github.com/usegro/services/crm/config"
	crmModels "github.com/usegro/services/crm/internal/apps/crm/models"
	"github.com/usegro/services/crm/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgressInstance1 *gorm.DB

func SetUpPostgres() {
	cfg := config.GetConfig().Postgres

	// Validate required config
	if cfg.Host == "" {
		logger.Log.Fatal("Postgres host is empty — skipping setup")
	}

	if cfg.Schema == "" {
		logger.Log.Fatal("Postgres schema is not set")
	}

	// Set up timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Connect using GORM
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=UTC",
		cfg.Host,     // host
		cfg.Username, // user
		cfg.Password, // password
		cfg.Database, // dbname
		cfg.Port,     // port
	)), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("gorm.Open() failed", zap.Error(err))
	}

	// Ping to confirm connection
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Fatal("db.DB() failed", zap.Error(err))
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Log.Fatal("Postgres ping failed", zap.Error(err))
	}

	logger.Log.Info("✅ GORM connection established")

	// AutoMigrate CRM-owned models only
	// Base models (User, Job, etc.) are migrated by the base service
	if err := db.AutoMigrate(
		&crmModels.CrmUserOrganization{},
		&crmModels.SalesChannelType{},
		&crmModels.StockProductType{},
	); err != nil {
		logger.Log.Fatal("AutoMigrate failed", zap.Error(err))
	}
	logger.Log.Info("✅ AutoMigration complete")

	PostgressInstance1 = db
}
