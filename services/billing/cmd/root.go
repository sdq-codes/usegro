package cmd

import (
	"fmt"
	"github.com/usegro/services/billing/database"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
	"github.com/usegro/services/billing/config"
	"github.com/usegro/services/billing/internal/logger"
	"go.uber.org/zap"
)

var RootCmdName = "main"

var Db *gorm.DB

var configFile string
var rootCmd = &cobra.Command{
	Use: func() string {
		return RootCmdName
	}(),
	Short: "\nThis application is made with ❤️",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Usage()
	},
}

func selectConfigFile() string {
	env := os.Getenv("APP_ENV")
	const defaultConfigFile = "config/config.dev.yaml"
	switch env {
	case "production", "prod":
		return filepath.FromSlash("config/config.prod.yaml")
	case "staging":
		return filepath.FromSlash("config/config.staging.yaml")
	case "development", "dev":
		return filepath.FromSlash("config/config.dev.yaml")
	default:
		fmt.Printf("⚠️  APP_ENV not set or unrecognized (%q). Using default: %s\n", env, defaultConfigFile)
		return filepath.FromSlash(defaultConfigFile)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is %s)", selectConfigFile()))
}

func setupAll() {
	setUpConfig()
	setUpLogger()
	setUpRedis()
	setUpSentry()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("rootCmd.Execute() Error: %v", err)
		os.Exit(1)
	}
}

func setUpConfig() {
	if configFile == "" {
		configFile = selectConfigFile()
	}

	log.Default().Printf("Using config file: %s", configFile)
	config.SetConfig(configFile)
}

func setUpLogger() {
	log.Default().Printf("Using log level: %s", config.GetConfig().Log.Level)
	logger.InitLogger("zap")
}

func setUpRedis() {
	// Create the database connection pool
	if config.GetConfig().Redis[0].Host != "" {
		logger.Log.Info("Initializing redis")
		log.Println(config.GetConfig().Redis)
		err := database.InitRedisClient(config.GetConfig().Redis)
		if err != nil {
			logger.Log.Fatal("rdb.InitRedisClient()", zap.Error(err))
		}
		logger.Log.Info("redis initialized")
	}
}

func setUpSentry() {
	// Don't initialize sentry if DSN is not set
	if config.GetConfig().Sentry.Dsn == "" {
		return
	}

	// Initialize sentry
	logger.Log.Info("Initializing Sentry " + config.GetConfig().Sentry.Dsn)
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.GetConfig().Sentry.Dsn,
		// BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
		// 	if hint.Context != nil {
		// 		if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
		// 			// You have access to the original Context if it panicked
		// 			fmt.Println(utils.CopyString(c.Hostname()))
		// 		}
		// 	}
		// 	fmt.Println(event)
		// 	return event
		// },
		Debug:            config.GetConfig().Sentry.Debug,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Environment:      config.GetConfig().Sentry.Environment,
		Release:          config.GetConfig().Sentry.Release,
	})

	if err != nil {
		logger.Log.Error("Create Sentry instant error: %v", zap.Error(err))
		return
	}

	logger.Log.Info("Create Sentry instant success")

	// send initial event to sentry with data
	sentry.CaptureMessage("Sentry initialized")

	defer sentry.Flush(2 * time.Second)
}
