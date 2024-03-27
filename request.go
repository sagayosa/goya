package goya

import (
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

	buildErr  []error
	buildURL  string
	buildBody []byte
}

func NewRequestBuilder(method, url string, opt *Option) *RequestBuider {
	return &RequestBuider{
		Method: method,
		URL:    url,
		Opt:    opt,
	}
}

func (b *RequestBuider) Build() *http.Request {
	if b.Opt.Json != nil {
		b.buildJson()
	}
	request, _ := http.NewRequest(b.Method, b.URL, nil)
	return request
}

func (b *RequestBuider) errHappen(err error) {
	b.buildErr = append(b.buildErr, err)
}

func (b *RequestBuider) buildJson() {
	bts, err := json.Marshal(b.Opt.Json)
	if err != nil {
		b.errHappen(err)
	}
	b.buildBody = bts
}

func (b *RequestBuider) buildParams() {
	var mp map[string]any
	tp := reflect.TypeOf(b.Opt.Params).Kind()
	if tp == reflect.Struct {
		mp = ConvertStructToMap(mp)
	} else if tp == reflect.Map {
		var ok bool
		mp, ok = b.Opt.Params.(map[string]any)
		if !ok {
			b.errHappen(fmt.Errorf("params is map but not the map[string]any"))
			return
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
		querys.Add(k, fmt.Sprintln(v))
	}

	parsedURL.RawQuery = querys.Encode()
	b.buildURL = parsedURL.String()
}
