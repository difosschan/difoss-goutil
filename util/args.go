package util

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Separator = '='

var args map[string]string

func init() {
	args = PickArgs()
}

// Use for pick up arguments into map, either `go run` or `go test`.
func PickArgs() map[string]string {
	args := os.Args
	if len(args) == 0 {
		args = flag.Args()
	}

	m := make(map[string]string)
	for _, arg := range args {
		index := strings.IndexByte(arg, Separator)
		if index != -1 {
			k := arg[:index]
			v := arg[index+1:]
			k = strings.TrimPrefix(k, "-")

			if _, existed := m[k]; ! existed {
				m[k] = v
			}
		}
	}
	return m
}

func GetFromArgs(key string, defaultValueOrNone ... string) (value string) {
	var ok bool
	if value, ok = args[key]; ! ok {
		if len(defaultValueOrNone) == 0 {
			panic(fmt.Sprintf(`parameter "%s" MUST be required, like %s=...`, key, key))
		}
		value = strings.Join(defaultValueOrNone, " ")
	}
	return
}

func GetBoolFromArgs(key string, defaultValueOrNone ...bool) (value bool) {
	var s string
	if len(defaultValueOrNone) == 0 {
		s = GetFromArgs(key)
	} else {
		s = GetFromArgs(key, strconv.FormatBool(defaultValueOrNone[0]))
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return
	}
	return b
}

func GetIntFromArgs(key string, defaultValueOrNone ...int64) (value int64) {
	var s string
	if len(defaultValueOrNone) == 0 {
		s = GetFromArgs(key)
	} else {
		s = GetFromArgs(key, strconv.FormatInt(defaultValueOrNone[0], 10))
	}
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	return i64
}

func GetUintFromArgs(key string, defaultValueOrNone ...uint64) (value uint64) {
	var s string
	if len(defaultValueOrNone) == 0 {
		s = GetFromArgs(key)
	} else {
		s = GetFromArgs(key, strconv.FormatUint(defaultValueOrNone[0], 10))
	}
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return
	}
	return u64
}
