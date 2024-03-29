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

// type OptionFunc func(opt *Option)

// WithJson will add data to the request body
// and set the Content-Type to application/json
// func WithJson(data any) OptionFunc {
// 	return func(opt *Option) {
// 		opt.Json = data
// 	}
// }

// WithParams will add params to the request URL
// params must be a structure or a map
// func WithParams(params any) OptionFunc {
// 	return func(opt *Option) {
// 		opt.Params = params
// 	}
// }

// WithForm will add data to the request body
// and set the Content-Type to multipart/form-data
// func WithForm(data any) OptionFunc {
// 	return func(opt *Option) {
// 		opt.FormData = data
// 	}
// }

type OptionFunc func(any) (BeforeBuildFunc, AfterBuildFunc)
type BeforeBuildFunc func(b *RequestBuider)
type AfterBuildFunc func(req *http.Request)

func WithJson(data any) OptionFunc {
	return func(a any) (BeforeBuildFunc, AfterBuildFunc) {
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

func WithForm(data any) OptionFunc {
	return func(a any) (BeforeBuildFunc, AfterBuildFunc) {
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

func WithParams(params any) OptionFunc {
	return func(a any) (BeforeBuildFunc, AfterBuildFunc) {
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
