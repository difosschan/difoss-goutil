package log

import (
	"github.com/difosschan/difoss-goutil/util/log/rotate"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Encoding string

const (
	CapitalColorful   Encoding = "capital_colorful"
	LowercaseColorful Encoding = "LowercaseColorful"
	Colorful          Encoding = "colorful" // Compared with CapitalColorful, it has same effect.
	Color             Encoding = "color"
	Json              Encoding = "json"
	Console           Encoding = "console"
)

type WriterConfig struct {
	Enable   bool          `json:"enable"`
	Level    zapcore.Level `json:"level"`
	Encoding Encoding      `json:"encoding"`
}

type Config struct {
	StdoutWriter    WriterConfig  `json:"stdout_writer,omitempty"`
	FileWriter      WriterConfig  `json:"file_writer,omitempty"`
	Rotate          rotate.Config `json:"rotate"`
	AddCaller       bool          `json:"add_caller"`
	StackTraceLevel zapcore.Level `json:"stack_trace_level"`
}

func DefaultConfig() *Config {
	return &Config{
		StdoutWriter: WriterConfig{
			Enable: true, Level: zap.DebugLevel, Encoding: Colorful},
		FileWriter: WriterConfig{
			Enable: false, Level: zap.InfoLevel, Encoding: Console},
		Rotate: *rotate.DefaultConfig(),

		StackTraceLevel: zap.WarnLevel,
	}
}

func (cfg *Config) Build(options ...zap.Option) *zap.Logger {
	var cores []zapcore.Core

	if cfg.StdoutWriter.Enable {
		cores = append(cores, zapcore.NewCore(
			getEncoder(cfg.StdoutWriter.Encoding), zapcore.AddSync(os.Stdout), cfg.StdoutWriter.Level))
	}

	if cfg.FileWriter.Enable {
		cores = append(cores, zapcore.NewCore(
			getEncoder(cfg.FileWriter.Encoding), rotate.GetLogWriter(&cfg.Rotate), cfg.FileWriter.Level))
	}
	core := zapcore.NewTee(cores...)

	if cfg.AddCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.StackTraceLevel > zapcore.InfoLevel {
		options = append(options, zap.AddStacktrace(cfg.StackTraceLevel))
	}
	logger := zap.New(core, options...)
	return logger
}

func getEncoder(encoding Encoding) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	var encoder zapcore.Encoder
	switch encoding {
	case Color, Colorful, CapitalColorful:
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case LowercaseColorful:
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case Json:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case Console:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}
