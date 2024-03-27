package goya

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

const (
	testURL = "http://httpbin.org/get"
)

type BasicGetResponse struct {
	Args    any     `json:"args"`
	Headers Headers `json:"headers"`
	Origin  string  `json:"origin"`
	URL     string  `json:"url"`
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

func TestDo(t *testing.T) {
	ts := []struct {
		url    string
		params any
		data   any
		want   *BasicGetResponse
	}{
		{
			testURL,
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

	for _, tt := range ts {
		resp := NewRequestClient("GET", tt.url, NewOption(WithJson(tt.data), WithParams(tt.params)), nil).Do()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Do got StatusCode %v but want %v", resp.StatusCode, http.StatusOK)
		}
		body, _ := io.ReadAll(resp.Body)
		basic := &BasicGetResponse{}
		json.Unmarshal(body, basic)

		if !compareResp(basic, tt.want) {
			t.Errorf("Do got body %v but want %v", basic, tt.want)
		}
	}
}
