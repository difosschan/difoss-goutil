package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MergeMode int

// Merge A -> D
const (
	// Overwrite the non-empty fields in A (addition) over the corresponding fields in D (dest).
	// For short: A(non-empty fields) ----> D
	OverwriteWithNonEmpty MergeMode = iota
	// Overwrite all fields owned by A (addition).
	// For short: A(all fields) ----> D
	Overwrite
	// Fill in the blank field of D (dest) with corresponding fields in A (addition).
	// In this way, fields that have non-empty values will not be affected.
	// For short: A(non-empty fields that is empty in D) ----> D
	FillBlank
)

var (
	TypeMustBeStruct = errors.New("type must be struct")
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

func MergeMap(dest, addition map[string]interface{}, mode MergeMode) {
	for k, v := range addition {
		addValue := reflect.ValueOf(v)

		if destValue, ok := dest[k]; ok {
			// log.Printf("IN dest[%s]=%v, reflect.ValueOf(destValue)=%v, destValue.IsZero=%v\n",
			// 	k, destValue, reflect.ValueOf(destValue), reflect.ValueOf(destValue).IsZero())

			if addValue.Kind() == reflect.Map {
				destMapValue, destOK := destValue.(map[string]interface{})
				addMapValue, addOK := v.(map[string]interface{})
				if destOK && addOK {
					// log.Printf("RECURSIVE CALL MergeMap(), k = %s \n", k)
					MergeMap(destMapValue, addMapValue, mode)
					dest[k] = destMapValue
					continue
				}
			}

			if mode == FillBlank {
				if ! reflect.ValueOf(destValue).IsZero() {
					continue
				}
			} else if mode == OverwriteWithNonEmpty {
				if addValue.IsZero() {
					continue
				}
			}
		}
		dest[k] = reflect.New(addValue.Type())
		dest[k] = v
	}
	return
}

func MergeStruct(out, in interface{}, mode MergeMode) (result interface{}, err error) {
	outType := reflect.TypeOf(out)
	outValue := reflect.ValueOf(out)

	if outValue.Kind() != reflect.Struct {
		return nil, TypeMustBeStruct
	}

	r := reflect.New(outType).Elem()

	outFieldCnt := outValue.NumField()
	for i := 0; i < outFieldCnt; i ++ {
		iOutName := outType.Field(i).Name
		iOutValue := outValue.Field(i)

		iResultValue := r.FieldByName(iOutName)
		iResultValue.Set(iOutValue)

		if !iOutValue.IsValid() {
			continue
		}

		iOutKind := iOutValue.Kind()
		iInValue := reflect.ValueOf(in).FieldByName(iOutName)
		// fmt.Printf("%d/%d: (kind=%v), out[%s] <<%s>> = %v\n",
		// 	i, outFieldCnt, iOutKind, iOutName, iOutValue.Type().String(), iOutValue)

		if !iInValue.IsValid() {
			continue // `in` did not contain field named by `iOutName`
		}
		// fmt.Printf("\t [%s] in  .Value: %v, .Kind(): %v\n",
		// 	iOutName, iInValue, iInValue.Kind())

		if iOutKind == reflect.Struct {
			if res, err := MergeStruct(iOutValue.Interface(), iInValue.Interface(), mode); err != nil {
				return nil, err
			} else {
				iResultValue.Set(reflect.ValueOf(res))
			}
		} else {
			if mode == FillBlank {
				if !iOutValue.IsZero() {
					continue
				}
			} else if mode == OverwriteWithNonEmpty {
				if iInValue.IsZero() {
					continue
				}
			}
			iResultValue.Set(iInValue)
		}
	}
	return r.Interface(), nil
}

func StructToString(s interface{}, showZero ...bool) (result string) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() != reflect.Struct {
		return
	}

	checkZero := len(showZero) == 0 || false == showZero[0]

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Name

		if strings.ToLower(f.Tag.Get("disable")) == "true" {
			continue
		}
		value := v.FieldByName(f.Name)
		if checkZero && value.IsZero() {
			continue
		}
		valueInterface := value.Interface()
		if value.Kind() == reflect.Struct {
			valueInterface = string("{") + StructToString(valueInterface, showZero...) + "}"
		}
		if i != 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v: %v", key, valueInterface)
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
