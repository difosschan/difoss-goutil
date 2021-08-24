package log

import (
	"encoding/json"
	"github.com/difosschan/difoss-goutil/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"testing"
)

var LogConfigYaml = `log:
  stdout_writer:
    enable: true
    encoding: colorful
    level: debug
  file_writer:
    enable: true
    encoding: json
    level: info
  rotate:
    filename: difoss-goutil.log
    is_compress: true
    is_local_time: true
    max_age_day: 10
    max_backup: 10
    max_size_mb: 1
`

func TestLoadConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigType("yaml")
	r := strings.NewReader(LogConfigYaml)
	if err := v.ReadConfig(r); err != nil {
		t.Fatal("Viper.ReadConfig failed:", err)
	}
	var i interface{}
	v.Unmarshal(&i)

	t.Log(util.JsonDump(i, 2))

	config := struct {
		Log Config
	}{}

	bs, _ := json.Marshal(&i)
	if err := json.Unmarshal(bs, &config); err != nil {
		t.Fatal("json.Unmarshal failed:", err)
	}

	log := config.Log.Build()
	if log != nil {
		zap.ReplaceGlobals(log)
	}
	defer log.Sync()

	zap.L().Debug("test debug log...")
	zap.L().Info("test info log...")
}
