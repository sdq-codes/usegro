package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var config *Config
var m sync.Mutex

type Config struct {
	Env        string     `yaml:"env"`
	App        App        `yaml:"apps"`
	FrontEnd   FrontEnd   `yaml:"frontEnd"`
	HttpServer HttpServer `yaml:"httpServer"`
	GrpcServer GrpcServer `yaml:"grpcServer"`
	Log        Log        `yaml:"log"`
	Scheduler  Scheduler  `yaml:"scheduler"`
	Schedules  []Schedule `yaml:"schedules"`
	Postgres   Postgres   `yaml:"postgres"`
	SSLMode    string     `yaml:"sslMode"`
	Redis      []Redis    `yaml:"redis"`
	Auth       Auth       `yaml:"auth"`
	Sentry     Sentry     `yaml:"sentry"`
	R2         R2Config   `yaml:"r2"`
}

type R2Config struct {
	AccountID       string `yaml:"accountId"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	BucketName      string `yaml:"bucketName"`
	PublicURL       string `yaml:"publicUrl"`
}

type GrpcServer struct {
	Port int `yaml:"port"`
}

type HttpServer struct {
	Port int `yaml:"port"`
}

type Log struct {
	Level           string `yaml:"level"`
	StacktraceLevel string `yaml:"stacktraceLevel"`
	FileEnabled     bool   `yaml:"fileEnabled"`
	FileSize        int    `yaml:"fileSize"`
	FilePath        string `yaml:"filePath"`
	FileCompress    bool   `yaml:"fileCompress"`
	MaxAge          int    `yaml:"maxAge"`
	MaxBackups      int    `yaml:"maxBackups"`
}

type Label struct {
	En string `json:"en"`
	Th string `json:"th"`
}

type App struct {
	Name     string `yaml:"name"`
	NameSlug string `yaml:"nameSlug"`
}

type FrontEnd struct {
	Url string `yaml:"url"`
}

type Postgres struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Schema          string `yaml:"schema"`
	SSLMode         string `yaml:"sslMode"`
	MaxConnections  int32  `yaml:"maxConnections"`
	MaxConnIdleTime int32  `yaml:"maxConnIdleTime"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type Auth struct {
	TokenExpiryMinutes     int    `yaml:"tokenExpiryMinutes"`
	RefreshTokenExpiryDays int    `yaml:"refreshTokenExpiryDays"`
	ApiSecret              string `yaml:"apiSecret"`
}

type Sentry struct {
	Dsn         string `yaml:"dsn"`
	Environment string `yaml:"environment"`
	Release     string `yaml:"release"`
	Debug       bool   `yaml:"debug"`
}

type Scheduler struct {
	Timezone string `yaml:"timezone"`
}

type Schedule struct {
	Job       string `yaml:"job"`
	Cron      string `yaml:"cron"`
	IsEnabled bool   `yaml:"isEnabled"`
}

type Authentication struct {
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func GetConfig() *Config {
	return config
}

func SetConfig(configFile string) {
	m.Lock()
	defer m.Unlock()

	// Load .env if present (silently ignored if missing)
	_ = godotenv.Load(".env")

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	// Allow env vars to override sensitive config values.
	// Env var names are listed explicitly to avoid accidental exposure.
	viper.AutomaticEnv()
	bindEnvs := map[string]string{
		"postgres.password":  "POSTGRES_PASSWORD",
		"postgres.sslMode":   "POSTGRES_SSL_MODE",
		"auth.apiSecret":     "AUTH_API_SECRET",
		"sentry.dsn":         "SENTRY_DSN",
		"r2.accountId":       "R2_ACCOUNT_ID",
		"r2.accessKeyId":     "R2_ACCESS_KEY_ID",
		"r2.secretAccessKey": "R2_SECRET_ACCESS_KEY",
		"r2.bucketName":      "R2_BUCKET_NAME",
		"r2.publicUrl":       "R2_PUBLIC_URL",
	}
	for key, env := range bindEnvs {
		if err := viper.BindEnv(key, env); err != nil {
			log.Fatalf("Failed to bind env var %s: %s", env, err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}
