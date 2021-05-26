package util

import (
	"fmt"
	"reflect"
	"strings"
)

func VectorToMap(vector interface{}, keys ...string) (mapResult *map[string]interface{}) {
	m := make(map[string]interface{})
	keyLen := len(keys)
	if keyLen == 0 {
		return nil
	}
	typVec := reflect.TypeOf(vector)
	kind := typVec.Kind()
	if ! (kind == reflect.Slice || kind == reflect.Array) {
		return nil
	}

	vecVal := reflect.ValueOf(vector)

	var realVector []reflect.Value
	for i := 0; i < vecVal.Len(); i++ {
		item := vecVal.Index(i)
		realVector = append(realVector, item)
	}
	fmt.Printf("realVector: %v\n", realVector)

	for _, item := range realVector {
		var fieldStrValues []string
		for _, k := range keys {
			fieldValue := item.FieldByName(k)
			if fieldValue.IsValid() {
				fieldStrValues = append(fieldStrValues, fieldValue.String())
			}
		}
		key := GenKeyForMap(fieldStrValues...)
		if _, ok := m[key]; ok {
			continue
		}
		m[key] = item.Interface()
	}
	return &m
}

func GenKeyForMap(values ...string) string {
	return strings.Join(values, "|")
}
