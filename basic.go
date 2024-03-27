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
