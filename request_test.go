package goya

import (
	"encoding/json"
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
