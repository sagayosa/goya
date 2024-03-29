package goya

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestDo(t *testing.T) {
	ts := []struct {
		url    string
		params any
		data   any
		want   *basicGetResponse
	}{
		{
			getURL,
			map[string]any{
				"host": "127.0.0.1",
				"db":   "test",
			},
			nil,
			&basicGetResponse{
				Args: map[string]any{
					"host": "127.0.0.1",
					"db":   "test",
				},
				URL: "http://httpbin.org/get?host=127.0.0.1&db=test",
				Headers: headers{
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
			&basicGetResponse{
				Args: map[string]any{
					"host":    "127.0.0.1",
					"db":      "test2",
					"version": "2",
				},
				URL: "http://httpbin.org/get?host=127.0.0.1&db=test2&version=2",
				Headers: headers{
					ContentType:   "application/json",
					ContentLength: "49",
				},
			},
		},
	}

	for _, tt := range ts {
		resp := NewRequestClient("GET", tt.url, NewOption(WithJson(tt.data), WithParams(tt.params)), nil).Do()
		if resp == nil {
			t.Fatal("Do got nil")
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Do got StatusCode %v but want %v", resp.StatusCode, http.StatusOK)
		}
		body, _ := io.ReadAll(resp.Body)
		basic := &basicGetResponse{}
		json.Unmarshal(body, basic)

		if !compareResp(basic, tt.want) {
			t.Errorf("Do got body %v but want %v", basic, tt.want)
		}
	}
}
