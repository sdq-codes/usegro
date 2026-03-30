package logger

import (
	"sync"

	"github.com/sdq-codes/usegro-api/config"
	sharedlogger "github.com/usegro/services/shared/pkg/logger"
	"go.uber.org/zap"
)

var Log *zap.Logger
var m sync.Mutex

func InitLogger(logDriver string) {
	m.Lock()
	defer m.Unlock()

	cfg := config.GetConfig().Log
	sharedlogger.InitLogger(sharedlogger.LogConfig{
		Level:           cfg.Level,
		StacktraceLevel: cfg.StacktraceLevel,
		FileEnabled:     cfg.FileEnabled,
		FileSize:        cfg.FileSize,
		FilePath:        cfg.FilePath,
		FileCompress:    cfg.FileCompress,
		MaxAge:          cfg.MaxAge,
		MaxBackups:      cfg.MaxBackups,
	})
	Log = sharedlogger.Log
}
