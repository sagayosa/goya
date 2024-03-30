package goya

import (
	"encoding/json"
	"net/http"
)

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

// // GetOpts send a get request to the URL with opt options,
// // The return value will be an instance of the type T you specified
// // Return value is a JSON deserialization of the response body of the request
// func GetOpts[T any](URL string, opt *Option) T {
// 	return Request[T](http.MethodGet, URL, opt)
// }

// // PostOpts send a post request to the URL with opt options,
// // The return value will be an instance of the type T you specified
// // Return value is a JSON deserialization of the response body of the request
// func PostOpts[T any](URL string, opt *Option) T {
// 	return Request[T](http.MethodPost, URL, opt)
// }

// // PutOpts send a put request to the URL with opt options,
// // The return value will be an instance of the type T you specified
// // Return value is a JSON deserialization of the response body of the request
// func PutOpts[T any](URL string, opt *Option) T {
// 	return Request[T](http.MethodPut, URL, opt)
// }

// // DeleteOpts send a delete request to the URL with opt options,
// // The return value will be an instance of the type T you specified
// // Return value is a JSON deserialization of the response body of the request
// func DeleteOpts[T any](URL string, opt *Option) T {
// 	return Request[T](http.MethodDelete, URL, opt)
// }

// // Get send a get request to the URL with params
// // params must be struct or map[string]any
// func Get[T any](URL string, params any) T {
// 	return GetOpts[T](URL, NewOption(WithParams(params)))
// }

// // Post send a post request to the URL with body in json
// // body must be struct or map[string]any
// func Post[T any](URL string, body any) T {
// 	return PostOpts[T](URL, NewOption(WithJson(body)))
// }

// // Put send a put request to the URL with body in json
// // body must be struct or map[string]any
// func Put[T any](URL string, body any) T {
// 	return PutOpts[T](URL, NewOption(WithJson(body)))
// }

// // Delete send a delete request to the URL with body in json
// // body must be struct or map[string]any
// func Delete[T any](URL string, body any) T {
// 	return DeleteOpts[T](URL, NewOption(WithJson(body)))
// }

// Get send a request to the URL
// If opt is *Option, the request will be built based on the specified options.
// If the opt is not *Option, the opt will be parsed as params
func Get[T any, R any | *Option](URL string, opt R) T {
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
func Post[T any, R any | *Option](URL string, opt R) T {
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
func Put[T any, R any | *Option](URL string, opt R) T {
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
func Delete[T any, R any | *Option](URL string, opt R) T {
	switch v := any(opt).(type) {
	case *Option:
		return Request[T](http.MethodDelete, URL, v)
	default:
		return Request[T](http.MethodDelete, URL, NewOption(WithJson(v)))
	}
}
