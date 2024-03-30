package goya

import (
	"encoding/json"
	"net/http"
)

// RequestRaw returns the *http.Response after the request and no reads were made
func RequestRaw(method, URL string, opt *Option) *Response {
	return NewRequestClient(method, URL, opt, nil).Do()
}

// Request returns the instance of the given T after JSON parsing
func Request[T any](method, URL string, opt *Option) T {
	bts, _ := NewRequestClient(method, URL, opt, nil).Do().Bytes()
	var result T
	json.Unmarshal(bts, &result)
	return result
}

// Get send a request to the URL
// If opt is *Option, the request will be built based on the specified options.
// If the opt is not *Option, the opt will be parsed as params
func Get[T any](URL string, opt any) T {
	switch v := any(opt).(type) {
	case *Option:
		return Request[T](http.MethodGet, URL, v)
	default:
		return Request[T](http.MethodGet, URL, NewOption(WithParams(v)))
	}
}

// Post send a request to the URL
// If opt is *Option, the request will be built based on the specified options.
// If the opt is not *Option, the opt will be parsed as json
func Post[T any](URL string, opt any) T {
	switch v := any(opt).(type) {
	case *Option:
		return Request[T](http.MethodPost, URL, v)
	default:
		return Request[T](http.MethodPost, URL, NewOption(WithJson(v)))
	}
}

// Put send a request to the URL
// If opt is *Option, the request will be built based on the specified options.
// If the opt is not *Option, the opt will be parsed as json
func Put[T any](URL string, opt any) T {
	switch v := any(opt).(type) {
	case *Option:
		return Request[T](http.MethodPut, URL, v)
	default:
		return Request[T](http.MethodPut, URL, NewOption(WithJson(v)))
	}
}

// Delete send a request to the URL
// If opt is *Option, the request will be built based on the specified options.
// If the opt is not *Option, the opt will be parsed as json
func Delete[T any](URL string, opt any) T {
	switch v := any(opt).(type) {
	case *Option:
		return Request[T](http.MethodDelete, URL, v)
	default:
		return Request[T](http.MethodDelete, URL, NewOption(WithJson(v)))
	}
}
