package database

import (
	"context"
	"fmt"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	formModels "github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
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

	//// Construct DSN (Data Source Name)
	//dsn := fmt.Sprintf(
	//	"host=%s port=%s authentication=%s password=%s dbname=%s sslmode=disable search_path=%s",
	//	cfg.Host,
	//	strconv.Itoa(cfg.Port),
	//	cfg.Username,
	//	cfg.Password,
	//	cfg.Database,
	//	cfg.Schema,
	//)

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

	//Optional ENUM SQL
	if err := models.CreateVerificationEnums(db); err != nil {
		return
	}
	logger.Log.Info("✅ ENUM SQL complete")

	// Optional: AutoMigrate models
	if err := db.AutoMigrate(&models.User{},
		&models.Job{},
		&models.VerificationToken{},
		&models.Verification{},
		&formModels.FieldType{},
		&formModels.FieldTypeConfig{},
		&formModels.FieldTypeValidation{}); err != nil {
		logger.Log.Fatal("AutoMigrate failed", zap.Error(err))
	}
	logger.Log.Info("✅ AutoMigration complete")

	// Optional: AutoMigrate models
	if err := AutoSeed(db); err != nil {
		logger.Log.Fatal("AutoMigrate failed", zap.Error(err))
	}
	logger.Log.Info("✅ AutoMigration complete")

	logger.Log.Info("✅ FieldTypes seeding complete")

	PostgressInstance1 = db
}

func AutoSeed(db *gorm.DB) error {
	if err := formModels.SeedFieldTypes(db); err != nil {
		logger.Log.Fatal("Seeding field types failed", zap.Error(err))
		return err
	}
	return nil
}
