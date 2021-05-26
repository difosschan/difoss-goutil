package util

import (
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var IsRelByBasedir = true

var baseDir string

type em struct{}

func init() {
	f, _, _, _ := GetFuncInfo(0)
	absDir, _ := filepath.Abs(f)
	curDir := path.Dir(absDir)
	curPkgPath := reflect.TypeOf(em{}).PkgPath()
	baseDir = strings.TrimSuffix(curDir, curPkgPath)
	baseDir = strings.TrimSuffix(baseDir, string(filepath.Separator))
}

func GetFuncInfo(skip int) (file string, line int, funcName string, ok bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if file != "" && IsRelByBasedir {
		file = RelativeDirByBase(file)
	}
	if ! ok {
		return
	}
	function := runtime.FuncForPC(pc)
	funcName = function.Name()
	return
}

func BaseDir() string {
	return baseDir
}

func RelativeDirByBase(targpath string) string {
	if strings.HasPrefix(targpath, baseDir) {
		s, _ := filepath.Rel(baseDir, targpath)
		return s
	}
	return targpath
}
