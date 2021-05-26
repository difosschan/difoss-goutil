package util

import (
	"encoding/json"
	"reflect"
	"strings"
)

func JsonDump(v interface{}, indent int) string {
	var b []byte
	var err error
	if indent > 0 {
		b, err = json.MarshalIndent(v, " ", strings.Repeat(" ", indent))
	} else {
		b, err = json.Marshal(v)
	}

	if err != nil {
		return ""
	}
	return string(b)
}

func StructureToMap(o interface{}) map[string]interface{} {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	maps := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		k := f.Name

		if f.Tag.Get("disable") == "true" {
			continue
		}
		v := v.FieldByName(f.Name).Interface()
		kind := reflect.TypeOf(v).Kind()
		// fmt.Println("==== v:", v, "v.Kind():", kind) // for DEBUG
		if kind == reflect.Struct {
			maps[k] = StructureToMap(v)
		} else {
			maps[k] = v
		}
	}
	return maps
}
