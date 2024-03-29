package goya

import (
	"bytes"
	"net/http"
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
	for _, before := range b.Opt.before {
		before(b)
	}

	request, err := http.NewRequest(b.method, b.URL, bytes.NewBuffer(b.Body))
	if err != nil {
		b.ErrHappen(err)
	}

	for _, after := range b.Opt.after {
		after(request)
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

// ErrHappen will add err to b.errs
func (b *RequestBuider) ErrHappen(err error) {
	b.errs = append(b.errs, err)
}
