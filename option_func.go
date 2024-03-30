package goya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type OptionFunc func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc)
type BeforeBuildFunc func(b *RequestBuider)
type AfterBuildFunc func(req *http.Request)
type ClientBuildFunc func(client *http.Client)
type ClientDoneFunc func(errs []error, resp *Response)

// WithJson will inject data into the body of the request in JSON format and set the Content-Type to application/json
// data can be struct or map
func WithJson(data any) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		if data == nil {
			return func(b *RequestBuider) { b.ErrHappen(fmt.Errorf("WithJson data is nil")) }, nil, nil, nil
		}
		return func(b *RequestBuider) {
				bts, err := json.Marshal(data)
				if err != nil {
					b.ErrHappen(err)
				}
				b.Body = bts
			}, func(req *http.Request) {
				req.Header.Set(contentType, contentTypeJSON)
			}, nil, nil
	}
}

// WithForm will inject data into the body of the request in form data and set the Content-Type to multipart/form-data
// data can be struct or map but the form data only support string and []string as values
// Therefore, if the value is not the string or []string, it will be changed to string by fmt.Sprintf() (may be JSON is better?)
func WithForm(data any) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		if data == nil {
			return func(b *RequestBuider) { b.ErrHappen(fmt.Errorf("WithForm data is nil")) }, nil, nil, nil
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		return func(b *RequestBuider) {
				mp, err := convertToMapStringAny(data)
				if err != nil {
					b.ErrHappen(err)
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
								b.ErrHappen(err)
							}
						}
					} else {
						err := writer.WriteField(k, fmt.Sprintf("%v", v))
						if err != nil {
							b.ErrHappen(err)
						}
					}
				}
				writer.Close()
				b.Body = body.Bytes()
			}, func(req *http.Request) {
				req.Header.Set(contentType, writer.FormDataContentType())
			}, nil, nil
	}
}

// WithParams will inject data into the URL as the query params
// data can be struct or map
// but the value will be changed to string by fmt.Sprintf() (may be JSON is better?)
func WithParams(params any) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		if params == nil {
			return func(b *RequestBuider) { b.ErrHappen(fmt.Errorf("WithParams params is nil")) }, nil, nil, nil
		}
		return func(b *RequestBuider) {
			mp, err := convertToMapStringAny(params)
			if err != nil {
				b.ErrHappen(err)
			}

			parsedURL, err := url.Parse(b.URL)
			if err != nil {
				b.ErrHappen(fmt.Errorf("URL is invalid : %w", err))
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
		}, nil, nil, nil
	}
}

// WithForceHeaders will inject all key value pairs from the headers into the *http.Request's Header
func WithForceHeaders(headers http.Header) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		return nil, func(req *http.Request) {
			for k, v := range headers {
				for i, h := range v {
					if i == 0 {
						req.Header.Set(k, h)
						continue
					}
					req.Header.Add(k, h)
				}
			}
		}, nil, nil
	}
}

// WithCookies will inject all cookies into the *http.Request's Cookie
func WithCookies(cookies []*http.Cookie) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		return nil, func(req *http.Request) {
			for _, c := range cookies {
				req.AddCookie(c)
			}
		}, nil, nil
	}
}

// WithTimeout will set timeout to *http.Client.Timeout
func WithTimeout(timeout time.Duration) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		return nil, nil, func(client *http.Client) {
			client.Timeout = timeout
		}, nil
	}
}

// WithForceHeader will set the header to *http.Request
func WithForceHeader(header string, value string) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		return nil, func(req *http.Request) {
			req.Header.Set(header, value)
		}, nil, nil
	}
}

// WithError will set the errors that occur during request generation and sending to errs
func WithError(errs *[]error) OptionFunc {
	return func() (BeforeBuildFunc, AfterBuildFunc, ClientBuildFunc, ClientDoneFunc) {
		return nil, nil, nil, func(errors []error, resp *Response) {
			*errs = errors
		}
	}
}
