package rotate

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Filename    string `json:"filename"`
	MaxSizeMB   int    `json:"max_size_mb"`
	MaxAgeDay   int    `json:"max_age_day"`
	MaxBackup   int    `json:"max_backup"`
	IsLocalTime bool   `json:"is_local_time"`
	IsCompress  bool   `json:"is_compress"`
}

func DefaultConfig() *Config {
	return &Config{
		Filename:    "default.log",
		MaxSizeMB:   1,
		MaxAgeDay:   30,
		MaxBackup:   7,
		IsCompress:  false,
		IsLocalTime: true,
	}
}

func GetLogWriter(cfg *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSizeMB,
		MaxAge:     cfg.MaxAgeDay,
		MaxBackups: cfg.MaxBackup,
		LocalTime:  cfg.IsLocalTime,
		Compress:   cfg.IsCompress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
