package goya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
		b.buildBody()
		b.buildURL()
	}

	request, err := http.NewRequest(b.Method, b.url, bytes.NewBuffer(b.body))
	if err != nil {
		b.errHappen(err)
	}

	if b.Opt != nil {
		b.buildHeaders(request)
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

func (b *RequestBuider) buildHeaders(req *http.Request) {
	for k, v := range b.Opt.Headers {
		req.Header.Set(k, v)
	}
}

func (b *RequestBuider) buildBody() {
	if b.Opt.Json != nil {
		b.buildJson()
	}
}

func (b *RequestBuider) buildURL() {
	if b.Opt.Params != nil {
		b.buildParams()
	}
}

func (b *RequestBuider) buildJson() {
	bts, err := json.Marshal(b.Opt.Json)
	if err != nil {
		b.errHappen(err)
	}
	b.body = bts
}

func (b *RequestBuider) buildFormData() {

}

func (b *RequestBuider) buildParams() {
	mp, err := ConvertToMapStringAny(b.Opt.Params)
	if err != nil {
		b.errHappen(err)
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
