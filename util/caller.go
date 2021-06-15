package util

import (
	"path/filepath"
	"runtime"
	"strings"
)

var IsRelByBasedir = true

var baseDir string

func init() {
	f, _, _, _ := GetFuncInfo(0)
	absDir, _ := filepath.Abs(f)
	baseDir = GetLinuxStylePath(absDir)
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
