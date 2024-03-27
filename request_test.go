package goya

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
)

func TestBuildJson(t *testing.T) {
	ts := []struct {
		url  string
		data any
	}{
		{
			"http://127.0.0.1:3306",
			map[string]any{
				"host": "127.0.0.1",
				"db":   "test",
			},
		},
		{
			"http://127.0.0.1:3306",
			struct {
				host    string
				db      string
				version int
			}{
				"127.0.0.1",
				"test2",
				2,
			},
		},
		{
			"http://127.0.0.1:3306",
			2,
		},
	}

	for _, tt := range ts {
		b := NewRequestBuilder("GET", tt.url, NewOption(WithJson(tt.data)))
		b.buildJson()

		want, _ := json.Marshal(tt.data)

		if !reflect.DeepEqual(b.buildBody, want) {
			t.Errorf("buildJson got %v but want %v", b.buildBody, want)
		}
	}
}

func TestBuildParams(t *testing.T) {
	ts := []struct {
		url  string
		data any
		want string
	}{
		{
			"http://127.0.0.1:3306",
			map[string]any{
				"host": "127.0.0.1",
				"db":   "test",
			},
			"http://127.0.0.1:3306?host=127.0.0.1&db=test",
		},
		{
			"http://127.0.0.1:3306",
			struct {
				Host    string `json:"host"`
				Db      string `json:"db"`
				Version int    `json:"version"`
			}{
				"127.0.0.1",
				"test2",
				2,
			},
			"http://127.0.0.1:3306?host=127.0.0.1&db=test2&version=2",
		},
		{
			"http://127.0.0.1:3306?temp=114514",
			struct {
				Host    string `json:"host"`
				Db      string `json:"db"`
				Version int    `json:"version"`
			}{
				"127.0.0.1",
				"test2",
				2,
			},
			"http://127.0.0.1:3306?host=127.0.0.1&db=test2&version=2&temp=114514",
		},
		{
			"http://127.0.0.1:3306",
			3306,
			"http://127.0.0.1:3306",
		},
	}

	for _, tt := range ts {
		b := NewRequestBuilder("GET", tt.url, NewOption(WithParams(tt.data)))
		b.buildParams()

		parsedWant, _ := url.Parse(tt.want)
		parsedGot, _ := url.Parse(b.buildURL)

		if !reflect.DeepEqual(parsedGot.Query(), parsedWant.Query()) {
			t.Errorf("buildParams got %v but want %v", b.buildURL, tt.want)
		}
	}
}
