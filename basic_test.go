package goya

import (
	"testing"
)

var ts = []struct {
	url    string
	params any
	data   any
	want   *BasicGetResponse
}{
	{
		getURL,
		map[string]any{
			"host": "127.0.0.1",
			"db":   "test",
		},
		nil,
		&BasicGetResponse{
			Args: map[string]any{
				"host": "127.0.0.1",
				"db":   "test",
			},
			URL: "http://httpbin.org/get?host=127.0.0.1&db=test",
			Headers: Headers{
				ContentType:   "",
				ContentLength: "",
			},
		},
	},
	{
		getURL,
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
		&BasicGetResponse{
			Args: map[string]any{
				"host":    "127.0.0.1",
				"db":      "test2",
				"version": "2",
			},
			URL: "http://httpbin.org/get?host=127.0.0.1&db=test2&version=2",
			Headers: Headers{
				ContentType:   "application/json",
				ContentLength: "49",
			},
		},
	},
}

func TestRequest(t *testing.T) {
	for _, tt := range ts {
		resp := Request[BasicGetResponse]("GET", tt.url, NewOption(WithJson(tt.data), WithParams(tt.params)))
		if !compareResp(&resp, tt.want) {
			t.Errorf("Do got body %v but want %v", resp, tt.want)
		}
	}
}
