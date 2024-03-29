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

type OptionFunc func() (BeforeBuildFunc, AfterBuildFunc)
type BeforeBuildFunc func(b *RequestBuider)
type AfterBuildFunc func(req *http.Request)

// WithJson will inject data into the body of the request in JSON format and set the Content-Type to application/json
// data can be struct or map
func WithJson(data any) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc) {
		if data == nil {
			return func(b *RequestBuider) { b.errHappen(fmt.Errorf("WithJson data is nil")) }, func(req *http.Request) {}
		}
		return func(b *RequestBuider) {
				bts, err := json.Marshal(data)
				if err != nil {
					b.errHappen(err)
				}
				b.Body = bts
			}, func(req *http.Request) {
				req.Header.Set(ContentType, ContentTypeJSON)
			}
	}
}

// WithForm will inject data into the body of the request in form data and set the Content-Type to multipart/form-data
// data can be struct or map but the form data only support string and []string as values
// Therefore, if the value is not the string or []string, it will be changed to string by fmt.Sprintf() (may be JSON is better?)
func WithForm(data any) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc) {
		if data == nil {
			return func(b *RequestBuider) { b.errHappen(fmt.Errorf("WithForm data is nil")) }, func(req *http.Request) {}
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		return func(b *RequestBuider) {
				mp, err := ConvertToMapStringAny(data)
				if err != nil {
					b.errHappen(err)
				}
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
			}, func(req *http.Request) {
				req.Header.Set(ContentType, writer.FormDataContentType())
			}
	}
}

// WithParams will inject data into the URL as the query params
// data can be struct or map
// but the value will be changed to string by fmt.Sprintf() (may be JSON is better?)
func WithParams(params any) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc) {
		if params == nil {
			return func(b *RequestBuider) { b.errHappen(fmt.Errorf("WithParams params is nil")) }, func(req *http.Request) {}
		}
		return func(b *RequestBuider) {
			mp, err := ConvertToMapStringAny(params)
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
		}, func(req *http.Request) {}
	}
}

func WithForceHeaders(headers http.Header) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc) {
		return func(b *RequestBuider) {}, func(req *http.Request) {
			for k, v := range headers {
				for i, h := range v {
					if i == 0 {
						req.Header.Set(k, h)
						continue
					}
					req.Header.Add(k, h)
				}
			}
		}
	}
}
