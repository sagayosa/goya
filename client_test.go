package goya

import (
	"net/http"
	"testing"
)

const (
	testURL = "http://httpbin.org"
)

func TestDo(t *testing.T) {
	ts := []struct {
		url    string
		params any
		data   any
	}{
		{
			testURL,
			map[string]any{
				"host": "127.0.0.1",
				"db":   "test",
			},
			nil,
		},
		{
			testURL,
			struct {
				Host    string `json:"host"`
				Db      string `json:"db"`
				Version int    `json:"version"`
			}{
				"127.0.0.1",
				"test2",
				2,
			},
			map[string]any{
				"select": "*",
				"from":   "db",
				"where":  "id = ?",
				"?":      3,
			},
		},
	}

	for _, tt := range ts {
		resp := NewRequestClient("GET", tt.url, NewOption(WithJson(tt.data), WithParams(tt.params)), nil).Do()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Do got StatusCode %v but want %v", resp.StatusCode, http.StatusOK)
		}
	}
}
