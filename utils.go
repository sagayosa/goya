package goya

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ConvertToMapStringAny is mainly used to convert the struct and map into map[string]any
// If src is neither a struct nor a map, it will return an error : params is neither struct nor map[string]any
func ConvertToMapStringAny(src any) (map[string]any, error) {
	result := map[string]any{}
	tp := reflect.TypeOf(src).Kind()
	if tp == reflect.Struct || tp == reflect.Pointer {
		result = ConvertStructToMap(src)
	} else if tp == reflect.Map {
		// var ok bool
		// mp, ok = b.Opt.Params.(map[string]any)
		// if !ok {
		// 	b.errHappen(fmt.Errorf("params is map but not the map[string]any"))
		// 	return
		// }
		rv := reflect.ValueOf(src)
		for _, k := range rv.MapKeys() {
			result[fmt.Sprintf("%v", k)] = rv.MapIndex(k)
		}
	} else {
		return nil, fmt.Errorf("params is neither struct nor map[string]any")
	}

	return result, nil
}

// src must be struct
// The index of the result will be the json in the tag of each field if the tag is exist
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
			if objVal.Type().Field(i).Tag.Get("json") != "" {
				result[objVal.Type().Field(i).Tag.Get("json")] = objVal.Field(i).Interface()
			} else {
				result[objVal.Type().Field(i).Name] = objVal.Field(i).Interface()
			}
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

type BasicGetResponse struct {
	Args    any     `json:"args"`
	Headers Headers `json:"headers"`
	Origin  string  `json:"origin"`
	URL     string  `json:"url"`
}

type BasicPostResponse struct {
	Headers Headers `json:"headers"`
	Origin  string  `json:"origin"`
	URL     string  `json:"url"`
	Args    any     `json:"args"`
	Data    string  `json:"data"`
	Files   any     `json:"files"`
	Form    any     `json:"form"`
}

type Headers struct {
	Accept         string `json:"Accept"`
	AcceptEncoding string `json:"Accept-Encoding"`
	AcceptLanguage string `json:"Accept-Language"`
	Host           string `json:"Host"`
	UserAgent      string `json:"User-Agent"`
	ContentType    string `json:"Content-Type"`
	ContentLength  string `json:"Content-Length"`
}

func compareResp(first *BasicGetResponse, second *BasicGetResponse) bool {
	if first.Headers.ContentType != second.Headers.ContentType {
		return false
	}
	if first.Headers.ContentLength != second.Headers.ContentLength {
		return false
	}
	parsedWant, _ := url.Parse(first.URL)
	parsedGot, _ := url.Parse(second.URL)
	if !reflect.DeepEqual(parsedWant.Query(), parsedGot.Query()) {
		return false
	}
	firstArgs, _ := json.Marshal(first.Args)
	secondArgs, _ := json.Marshal(second.Args)

	return reflect.DeepEqual(firstArgs, secondArgs)
}
