package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"strconv"
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

func MergeMap(dest, addition map[string]interface{}) {
	for k, v := range addition {
		value := reflect.ValueOf(v)
		isZero := value.IsZero()
		if !isZero {
			if _, ok := dest[k]; ok {
				dest[k] = reflect.New(value.Type())
				dest[k] = v
			}
		}
	}
	return
}

func MergeStruct(dest, addition interface{}) error {
	t := reflect.TypeOf(addition)
	v := reflect.ValueOf(addition)

	if t.Kind() != reflect.Struct {
		return errors.New("Type must be struct")
	}

	for i := 0; i < t.NumField(); i ++ {
		f := t.Field(i)

		iValue := v.Field(i)
		kind := iValue.Kind()
		zap.S().Debug("v:", v, "v.Kind():", kind)
		if kind == reflect.Struct {
			if e := MergeStruct(
				reflect.ValueOf(dest).FieldByName(f.Name),
				StructToMap(iValue)); e != nil {
				return e
			}
		} else {
			if !iValue.IsNil() && !iValue.IsZero() {
				reflect.ValueOf(dest).FieldByName(f.Name).Set(iValue)
			}
		}

	}
	return nil
}

func StructToString(s interface{}) (result string) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Name

		if strings.ToLower(f.Tag.Get("disable")) == "true" {
			continue
		}
		value := v.FieldByName(f.Name).Interface()
		if i != 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v: %v", key, value)
	}
	return
}

func StructToMap(o interface{}) map[string]interface{} {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	maps := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Name

		value := v.Field(i)

		// fmt.Println("==== key: ", key, "value:", value, "v.Kind():", value.Kind())

		if value.Kind() == reflect.Struct {
			maps[key] = StructToMap(value.Interface())
		} else if value.Kind() == reflect.Array {
			var valueMapFromArr = make(map[string]interface{})
			vNum := value.NumField()
			for iArr := 0; iArr < vNum; iArr ++ {
				valueMapFromArr[strconv.Itoa(iArr)] = value.Field(iArr)
			}
			maps[key] = StructToMap(valueMapFromArr)
		} else {
			maps[key] = value.Interface()
		}
	}
	return maps
}
