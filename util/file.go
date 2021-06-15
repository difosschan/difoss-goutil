package util

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrPathUnauthorized = errors.New("access path unauthorized")
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	return err == nil || os.IsExist(err)
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	return err == nil && ! s.IsDir()
}

func GetFileInfo(path string) (abs, dir, pureFilename, suffix string) {
	abs, _ = filepath.Abs(path)
	dir = filepath.Dir(abs)
	suffix = filepath.Ext(abs)
	if suffix != "" {
		pureFilename = strings.TrimSuffix(filepath.Base(abs), suffix)
	} else {
		pureFilename = filepath.Base(abs)
	}
	return
}

func FilePathJoinSafely(base string, children ...string) (string, error) {
	basePath, e := filepath.Abs(base)
	if e != nil {
		return "", e
	}
	childDir := filepath.Join(children...)
	joinedDir, e := filepath.Abs(filepath.Join(basePath, childDir))
	if e != nil {
		return "", e
	}

	if ! strings.HasPrefix(joinedDir, basePath) {
		return "", ErrPathUnauthorized
	}
	return joinedDir, nil
}

func GetLinuxStylePath(path string) (linuxPath string) {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, "\\", "/")
	}
	return path
}
