package goya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
)

type RequestBuider struct {
	method    string
	originURL string
	Opt       *Option

	errs []error
	// URL will be passed to NewRequest to create *http.Request,
	// so you can directly modify this field to get expected request
	URL string
	// Body will be passed to NewRequest to create *http.Request,
	// so you can directly modify this field to get expected request
	Body []byte
}

func NewRequestBuilder(method, url string, opt *Option) *RequestBuider {
	return &RequestBuider{
		method:    method,
		originURL: url,
		Opt:       opt,

		errs: make([]error, 0),
		URL:  url,
	}
}

func (b *RequestBuider) Build() *http.Request {
	if b.Opt != nil {
		b.buildBody()
		b.buildURL()
	}
	request, err := http.NewRequest(b.method, b.URL, bytes.NewBuffer(b.Body))
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
	} else if b.Opt.FormData != nil {
		b.buildFormData()
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
	b.Body = bts

	b.Opt.Headers[ContentType] = ContentTypeJSON
}

func (b *RequestBuider) buildFormData() {
	mp, err := ConvertToMapStringAny(b.Opt.FormData)
	if err != nil {
		b.errHappen(err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// The form data must be of type map[string]string or map[string][]string.
	// Therefore, I distinguish between the slice and other types, and use fmt.Sprintf directly.
	for k, v := range mp {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				element := val.Index(i)
				err := writer.WriteField(k, fmt.Sprintf("%v", element))
				if err != nil {
					b.errHappen(err)
				}
			}
		} else {
			err := writer.WriteField(k, fmt.Sprintf("%v", v))
			if err != nil {
				b.errHappen(err)
			}
		}
	}
	writer.Close()
	b.Body = body.Bytes()

	b.Opt.Headers[ContentType] = writer.FormDataContentType()
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
	// The params must be string
	// Therefore, i use fmt.Sprintf directly
	for k, v := range mp {
		querys.Add(k, fmt.Sprintf("%v", v))
	}

	parsedURL.RawQuery = querys.Encode()
	b.URL = parsedURL.String()
}
