package goya

import (
	"reflect"
	"strings"
)

// src must be struct
// The index of the result will be the json in the tag of each field
func ConvertStructToMap(src any) map[string]any {
	if src == nil {
		return nil
	}

	result := make(map[string]any)
	objVal := reflect.ValueOf(src)

	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}
	if objVal.Kind() == reflect.Struct {
		for i := 0; i < objVal.NumField(); i++ {
			result[objVal.Type().Field(i).Tag.Get("json")] = objVal.Field(i).Interface()
		}

		return result
	}

	return nil
}

func StringPlus(strs ...string) string {
	var builder strings.Builder
	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}
