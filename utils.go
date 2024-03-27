package goya

import (
	"reflect"
)

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
			result[objVal.Type().Field(i).Name] = objVal.Field(i).Interface()
		}

		return result
	}

	return nil
}
