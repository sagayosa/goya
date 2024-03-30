package goya

import (
	"net/http"
)

type RequestClient struct {
	Method  string
	URL     string
	Opt     *Option
	Client  *http.Client
	Request *http.Request

	errs []error
}

func NewRequestClient(method, url string, opt *Option, client *http.Client) *RequestClient {
	return &RequestClient{
		Method:  method,
		URL:     url,
		Opt:     opt,
		Client:  client,
		Request: nil,

		errs: make([]error, 0),
	}
}

// BuildRequest will build the request according to the Option and return the built http request
// Your modifications to the return value will be reflected in the request
func (c *RequestClient) BuildRequest() *http.Request {
	builder := NewRequestBuilder(c.Method, c.URL, c.Opt)
	request := builder.Build()
	if builder.Errors() != nil {
		c.errs = append(c.errs, builder.Errors()...)
	}
	c.Request = request
	return c.Request
}

// BuildClient will build the client according to the Option and return the built http client
// Your modifications to the return value will be reflected in the client
func (c *RequestClient) BuildClient() *http.Client {
	client := &http.Client{}
	for _, f := range c.Opt.client {
		f(client)
	}
	c.Client = client
	return c.Client
}

func (c *RequestClient) Do() *Response {
	if c.Request == nil {
		c.BuildRequest()
	}
	if c.Client == nil {
		c.BuildClient()
	}
	resp, err := c.Client.Do(c.Request)
	if err != nil {
		c.ErrHappen(err)
	}
	return NewResponse(resp)
}

// Return all errors that occurred during the Do()
// If no error occurs, return nil
func (c *RequestClient) Errors() []error {
	if len(c.errs) == 0 {
		return nil
	}
	return c.errs
}

// ErrHappen will add err to c.errs
func (c *RequestClient) ErrHappen(err error) {
	c.errs = append(c.errs, err)
}
