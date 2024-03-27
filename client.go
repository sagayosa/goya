package goya

import "net/http"

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

// Build will build the request according to the Option and return the built http request
// Your modifications to the return value will be reflected in the client
func (c *RequestClient) Build() *http.Request {
	request := NewRequestBuilder(c.Method, c.URL, c.Opt).Build()
	c.Request = request
	return c.Request
}

func (c *RequestClient) Do() *http.Response {
	if c.Request == nil {
		c.Build()
	}
	if c.Client == nil {
		c.Client = &http.Client{}
	}
	resp, err := c.Client.Do(c.Request)
	if err != nil {
		c.errHappen(err)
	}
	return resp
}

// Return all errors that occurred during the Do()
// If no error occurs, return nil
func (c *RequestClient) Errors() []error {
	if len(c.errs) == 0 {
		return nil
	}
	return c.errs
}

func (c *RequestClient) errHappen(err error) {
	c.errs = append(c.errs, err)
}
