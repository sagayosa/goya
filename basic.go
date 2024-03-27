package goya

import "net/http"

func RequestRaw(method, URL string, opt *Option) *http.Response {
	return NewRequestClient(method, URL, opt, nil).Do()
}
