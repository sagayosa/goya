package goya

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestWithJson(t *testing.T) {
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
				Host    string `json:"host"`
				Db      string `json:"db"`
				Version int    `json:"version"`
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
		b.Build()

		want, _ := json.Marshal(tt.data)

		if !reflect.DeepEqual(b.Body, want) {
			t.Errorf("buildJson got %v but want %v", b.Body, want)
		}
	}
}

func TestWithParams(t *testing.T) {
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
		b.Build()

		parsedWant, _ := url.Parse(tt.want)
		parsedGot, _ := url.Parse(b.URL)

		if !reflect.DeepEqual(parsedGot.Query(), parsedWant.Query()) {
			t.Errorf("Build got %v but want %v", b.URL, tt.want)
		}
	}
}

func TestWithFormData(t *testing.T) {
	ts := []struct {
		url  string
		data any
	}{
		{
			"http://127.0.0.1:3306",
			map[string]any{
				"host": "127.0.0.1",
				"db":   "test",
				"form": []string{"1", "2", "goya!!!!!"},
			},
		},
		{
			"http://127.0.0.1:3306",
			struct {
				Host    string `json:"host"`
				Db      string `json:"db"`
				Version string `json:"version"`
			}{
				"127.0.0.1",
				"test2",
				"2",
			},
		},
	}

	for _, tt := range ts {
		b := NewRequestBuilder("POST", tt.url, NewOption(WithForm(tt.data)))
		// b.buildFormData()
		r := b.Build()
		if err := r.ParseMultipartForm(20 << 32); err != nil {
			t.Error(err)
		}

		want, err := ConvertToMapStringAny(tt.data)
		if err != nil {
			t.Error(err)
		}
		form := ConvertFormToNormalOne(r.Form)

		if !reflect.DeepEqual(form, want) {
			t.Errorf("Build got %v but want %v", form, want)
		}
	}
}

func TestWithForceHeaders(t *testing.T) {
	ts := []struct {
		url     string
		headers map[string][]string
	}{
		{
			"http://127.0.0.1:3306",
			map[string][]string{"T": {"R"}, "R": {"T"}, "RT": {"T", "R"}, "TR": {"R", "T"}},
		},
	}
	for _, tt := range ts {
		b := NewRequestBuilder("GET", tt.url, NewOption(WithForceHeaders(tt.headers)))
		r := b.Build()

		for k, v := range tt.headers {
			rv, ok := r.Header[http.CanonicalHeaderKey(k)]
			if !ok {
				t.Errorf("header %v not set", k)
			}
			if !reflect.DeepEqual(v, rv) {
				t.Errorf("header %v got %v but want %v", k, rv, v)
			}
		}
	}
}

func TestWithCookies(t *testing.T) {
	ts := []struct {
		url     string
		cookies []*http.Cookie
	}{
		{
			"http://127.0.0.1:3306",
			[]*http.Cookie{{Name: "Test", Value: "123"}, {Name: "Test2", Value: "321"}},
		},
	}

	for _, tt := range ts {
		b := NewRequestBuilder("GET", tt.url, NewOption(WithCookies(tt.cookies)))
		r := b.Build()

		cookies := r.Cookies()
		if !reflect.DeepEqual(cookies, tt.cookies) {
			t.Errorf("cookies got %v but want %v", cookies, tt.cookies)
		}
	}
}
