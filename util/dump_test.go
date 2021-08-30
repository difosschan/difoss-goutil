package util

import (
	"github.com/difosschan/difoss-goutil/util/log"
	"github.com/difosschan/difoss-goutil/util/log/rotate"
	"github.com/magiconair/properties/assert"
	"go.uber.org/zap"
	"testing"
)

func TestMergeMapOverwriteWithNonEmpty(t *testing.T) {
	var a = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"enable":   true,
				"encoding": "colorful",
			},
			"rotate": map[string]interface{}{
				"is_compress": true,
				"max_age_day": 10,
			},
		},
	}

	var b = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"level":    "info", // will add
				"encoding": "json", // will modify
			},
			"rotate": map[string]interface{}{
				"is_compress": false, // no effect cause of empty value
				"max_age_day": 5,     // will modify
			},
			"addition_key": 123, // will add
			"addition_map": map[string]interface{}{ // will add
				"i": 1000,
				"s": "string",
				"map": map[string]interface{}{
					"d": 0.2,
				},
			},
		},
	}

	var target = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"enable":   true,
				"level":    "info",
				"encoding": "json",
			},
			"rotate": map[string]interface{}{
				"is_compress": true,
				"max_age_day": 5,
			},
			"addition_key": 123, // will add
			"addition_map": map[string]interface{}{ // will add
				"i": 1000,
				"s": "string",
				"map": map[string]interface{}{
					"d": 0.2,
				},
			},
		},
	}

	t.Log("a =", JsonDump(a, 2))
	t.Log("b =", JsonDump(b, 2))
	MergeMap(a, b, OverwriteWithNonEmpty)
	t.Log("AFTER MergeMap(how: OverwriteWithNonEmpty), a =", JsonDump(a, 2))
	assert.Equal(t, a, target)
}

func TestMergeMapOverwrite(t *testing.T) {
	var a = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"enable":   true,
				"encoding": "colorful",
			},
			"rotate": map[string]interface{}{
				"is_compress": true,
				"max_age_day": 10,
			},
		},
	}

	var b = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"level":    "info", // will add
				"encoding": "json", // will modify
			},
			"rotate": map[string]interface{}{
				"is_compress": false, // zero value but still be modified
				"max_age_day": 0,     // zero value but still be modified
			},
			"addition_key": 123, // will add
			"addition_map": map[string]interface{}{ // will add
				"i": 1000,
				"s": "string",
				"map": map[string]interface{}{
					"d": 0.2,
				},
			},
		},
	}

	var target = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"enable":   true,
				"level":    "info",
				"encoding": "json",
			},
			"rotate": map[string]interface{}{
				"is_compress": false,
				"max_age_day": 0,
			},
			"addition_key": 123, // will add
			"addition_map": map[string]interface{}{ // will add
				"i": 1000,
				"s": "string",
				"map": map[string]interface{}{
					"d": 0.2,
				},
			},
		},
	}

	t.Log("a =", JsonDump(a, 2))
	t.Log("b =", JsonDump(b, 2))
	MergeMap(a, b, Overwrite)
	t.Log("AFTER MergeMap(how: OverwriteWithNonEmpty), a =", JsonDump(a, 2))
	assert.Equal(t, a, target)
}

func TestMergeMapFillBlank(t *testing.T) {
	var a = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"enable":   true,
				"encoding": "colorful",
			},
			"rotate": map[string]interface{}{
				"is_compress": true,
				"max_age_day": 10,
			},
		},
	}

	var b = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"level":    "info", // will add
				"encoding": "json", // no effect cause of non-empty field in dest
			},
			"rotate": map[string]interface{}{
				"is_compress": false, // no effect cause of empty value
				"max_age_day": 5,     // no effect cause of non-empty field in dest
			},
			"addition_key": 123, // will add
			"addition_map": map[string]interface{}{ // will add
				"i": 1000,
				"s": "string",
				"map": map[string]interface{}{
					"d": 0.2,
				},
			},
		},
	}

	var target = map[string]interface{}{
		"log": map[string]interface{}{
			"stdout_writer": map[string]interface{}{
				"enable":   true,
				"level":    "info",
				"encoding": "colorful",
			},
			"rotate": map[string]interface{}{
				"is_compress": true,
				"max_age_day": 10,
			},
			"addition_key": 123,
			"addition_map": map[string]interface{}{
				"i": 1000,
				"s": "string",
				"map": map[string]interface{}{
					"d": 0.2,
				},
			},
		},
	}

	t.Log("a =", JsonDump(a, 2))
	t.Log("b =", JsonDump(b, 2))
	MergeMap(a, b, FillBlank)
	t.Log("AFTER MergeMap(how: FillBlank), a =", JsonDump(a, 2))
	assert.Equal(t, a, target)
}

func TestMergeStructOverwriteWithNonEmpty(t *testing.T) {
	var a = log.Config{
		StdoutWriter: log.WriterConfig{true, zap.DebugLevel, log.Colorful},
		FileWriter:   log.WriterConfig{false, zap.InfoLevel, log.Console},
		Rotate: rotate.Config{
			Filename:    "default.log",
			MaxSizeMB:   1,
			MaxAgeDay:   30,
			MaxBackup:   7,
			IsCompress:  false,
			IsLocalTime: true,
		},
	}
	var b = log.Config{
		FileWriter: log.WriterConfig{
			Enable:   true,
			Encoding: log.Json,
		},
		Rotate: rotate.Config{
			Filename:   "test.log",
			IsCompress: true,
			MaxAgeDay:  10,
		},
	}
	var target = log.Config{
		StdoutWriter: log.WriterConfig{
			Enable:   true,
			Level:    zap.DebugLevel,
			Encoding: log.Colorful,
		},
		FileWriter: log.WriterConfig{
			Enable:   true,
			Level:    zap.InfoLevel,
			Encoding: log.Json,
		},
		Rotate: rotate.Config{
			Filename:    "test.log",
			MaxSizeMB:   1,
			MaxAgeDay:   10,
			MaxBackup:   7,
			IsCompress:  true,
			IsLocalTime: true,
		},
	}
	t.Log("a =", JsonDump(a, 2))
	t.Log("b =", JsonDump(b, 2))
	c, err := MergeStruct(a, b, OverwriteWithNonEmpty)
	if err != nil {
		t.Fatal("MergeStruct fail:", err)
	}
	t.Log("After MergeStruct(how: OverwriteWithNonEmpty), c =", JsonDump(c, 2))
	assert.Equal(t, c, target)
}

func TestMergeStructFillBlank(t *testing.T) {
	var a = log.Config{
		StdoutWriter: log.WriterConfig{true, zap.DebugLevel, log.Colorful},
		FileWriter:   log.WriterConfig{false, zap.InfoLevel, log.Console},
		Rotate: rotate.Config{
			Filename:    "default.log",
			MaxSizeMB:   1,
			MaxAgeDay:   30,
			MaxBackup:   7,
			IsCompress:  false,
			IsLocalTime: true,
		},
	}
	var b = log.Config{
		FileWriter: log.WriterConfig{
			Enable:   true,
			Encoding: log.Json,
		},
		Rotate: rotate.Config{
			Filename:   "test.log",
			IsCompress: true,
			MaxAgeDay:  10,
		},
	}
	var target = log.Config{
		StdoutWriter: log.WriterConfig{
			Enable:   true,
			Level:    zap.DebugLevel,
			Encoding: log.Colorful,
		},
		FileWriter: log.WriterConfig{
			Enable:   true,
			Level:    zap.InfoLevel,
			Encoding: log.Console,
		},
		Rotate: rotate.Config{
			Filename:    "default.log",
			MaxSizeMB:   1,
			MaxAgeDay:   30,
			MaxBackup:   7,
			IsCompress:  true,
			IsLocalTime: true,
		},
	}
	t.Log("a =", JsonDump(a, 2))
	t.Log("b =", JsonDump(b, 2))
	c, err := MergeStruct(a, b, FillBlank)
	if err != nil {
		t.Fatal("MergeStruct fail:", err)
	}
	t.Log("After MergeStruct(how: FillBlank), c =", JsonDump(c, 2))
	assert.Equal(t, c, target)
}
