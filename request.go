package goya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type RequestBuider struct {
	Method string
	URL    string
	Opt    *Option

	errs []error
	url  string
	body []byte
}

func NewRequestBuilder(method, url string, opt *Option) *RequestBuider {
	return &RequestBuider{
		Method: method,
		URL:    url,
		Opt:    opt,

		errs: make([]error, 0),
		url:  url,
	}
}

func (b *RequestBuider) Build() *http.Request {
	if b.Opt != nil {
		if b.Opt.Json != nil {
			b.buildJson()
		}
		if b.Opt.Params != nil {
			b.buildParams()
		}
	}

	request, err := http.NewRequest(b.Method, b.url, bytes.NewBuffer(b.body))
	if err != nil {
		b.errHappen(err)
	}

	if b.Opt != nil {
		if len(b.Opt.Headers) != 0 {
			b.buildHeaders(request)
		}
	}
	return request
}

// Return all errors that occurred during the Build()
// If no error occurs, return nil
func (b *RequestBuider) Errors() []error {
	if len(b.errs) == 0 {
		return nil
	}
	return b.errs
}

func (b *RequestBuider) errHappen(err error) {
	b.errs = append(b.errs, err)
}

func (b *RequestBuider) buildJson() {
	bts, err := json.Marshal(b.Opt.Json)
	if err != nil {
		b.errHappen(err)
	}
	b.body = bts
}

func (b *RequestBuider) buildParams() {
	mp := map[string]any{}
	tp := reflect.TypeOf(b.Opt.Params).Kind()
	if tp == reflect.Struct || tp == reflect.Pointer {
		mp = ConvertStructToMap(b.Opt.Params)
	} else if tp == reflect.Map {
		// var ok bool
		// mp, ok = b.Opt.Params.(map[string]any)
		// if !ok {
		// 	b.errHappen(fmt.Errorf("params is map but not the map[string]any"))
		// 	return
		// }
		rv := reflect.ValueOf(b.Opt.Params)
		for _, k := range rv.MapKeys() {
			mp[fmt.Sprintf("%v", k)] = rv.MapIndex(k)
		}
	} else {
		b.errHappen(fmt.Errorf("params is neither struct nor map[string]any"))
		return
	}

	parsedURL, err := url.Parse(b.URL)
	if err != nil {
		b.errHappen(fmt.Errorf("URL is invalid : %w", err))
		return
	}
	querys := parsedURL.Query()
	for k, v := range mp {
		querys.Add(k, fmt.Sprintf("%v", v))
	}

	parsedURL.RawQuery = querys.Encode()
	b.url = parsedURL.String()
}

func (b *RequestBuider) buildHeaders(req *http.Request) {
	for k, v := range b.Opt.Headers {
		req.Header.Set(k, v)
	}
}
