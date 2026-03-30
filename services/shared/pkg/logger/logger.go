package logger

import (
	"sync"

	"go.uber.org/zap"
)

type LogConfig struct {
	Level           string
	StacktraceLevel string
	FileEnabled     bool
	FileSize        int
	FilePath        string
	FileCompress    bool
	MaxAge          int
	MaxBackups      int
}

var Log *zap.Logger
var m sync.Mutex

func InitLogger(cfg LogConfig) {
	m.Lock()
	defer m.Unlock()

	Log = newZapLogger(cfg)
}
