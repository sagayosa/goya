package goya

import "encoding/json"

// RequestRaw returns the *http.Response after the request and no reads were made
func RequestRaw(method, URL string, opt *Option) *Response {
	return NewResponse(NewRequestClient(method, URL, opt, nil).Do())
}

// Request returns the instance of the given T after JSON parsing
func Request[T any](method, URL string, opt *Option) T {
	bts, _ := NewResponse(NewRequestClient(method, URL, opt, nil).Do()).Bytes()
	var result T
	json.Unmarshal(bts, &result)
	return result
}

func Get[T any](URL string, opts ...OptionFunc) T {
	opt := NewOption(opts...)
	return Request[T]("GET", URL, opt)
}

func Post[T any](URL string, opts ...OptionFunc) T {
	opt := NewOption(opts...)
	return Request[T]("Post", URL, opt)
}

func Put[T any](URL string, opts ...OptionFunc) T {
	opt := NewOption(opts...)
	return Request[T]("Put", URL, opt)
}

func Delete[T any](URL string, opts ...OptionFunc) T {
	opt := NewOption(opts...)
	return Request[T]("Delete", URL, opt)
}
