package log

import (
	"github.com/difosschan/difoss-goutil/util/log/rotate"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type WriterConfig struct {
	Enable   bool   `json:"enable"`
	Level    string `json:"level"` // 可选值: debug（默认）, info, warn, error, dpanic, panic, fatal
	Encoding string `json:"encoding"`
}

type Config struct {
	StdoutWriter WriterConfig  `json:"stdout_writer,omitempty"`
	FileWriter   WriterConfig  `json:"file_writer,omitempty"`
	Rotate       rotate.Config `json:"rotate"`
}

func DefaultConfig() *Config {
	return &Config{
		StdoutWriter: WriterConfig{true, "debug", "colorful"},
		FileWriter:   WriterConfig{false, "info", "console"},
		Rotate:       *rotate.DefaultConfig(),
	}
}

func unmarshalLevel(level string) zapcore.Level {
	var e error
	var zapLevel zapcore.Level
	if e = zapLevel.Set(level); e != nil {
		zapLevel = zapcore.DebugLevel
	}
	return zapLevel
}

func (cfg *Config) Build() *zap.Logger {
	var cores []zapcore.Core

	if cfg.StdoutWriter.Enable {
		cores = append(cores, zapcore.NewCore(
			getEncoder(cfg.StdoutWriter.Encoding), zapcore.AddSync(os.Stdout), unmarshalLevel(cfg.StdoutWriter.Level)))
	}

	if cfg.FileWriter.Enable {
		cores = append(cores, zapcore.NewCore(
			getEncoder(cfg.FileWriter.Encoding), rotate.GetLogWriter(&cfg.Rotate), unmarshalLevel(cfg.FileWriter.Level)))
	}
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func getEncoder(encoding string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	var encoder zapcore.Encoder
	switch encoding {
	case "colorful":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}
