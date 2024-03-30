package goya

import (
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	// StatusCode is taken from *http.Response.StatusCode
	StatusCode int
	// Header is taken from *http.Response.Header
	Header      http.Header
	RawResponse *http.Response

	// You can get it after using Bytes() or String()
	Body []byte
}

// Bytes will read the body and return the result in []byte
// and the RawResponse will be closed
func (r *Response) Bytes() ([]byte, error) {
	if r.Body != nil {
		return r.Body, nil
	}
	if r.RawResponse == nil {
		return nil, fmt.Errorf("RawResponse is nil")
	}
	bts, err := io.ReadAll(r.RawResponse.Body)
	if err != nil {
		return bts, nil
	}
	r.RawResponse.Body.Close()
	r.Body = bts
	return r.Body, nil
}

// String will read the body and return the result in string
// and the RawResponse will be closed
func (r *Response) String() (string, error) {
	if r.Body != nil {
		return string(r.Body), nil
	}
	_, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(r.Body), nil
}

func NewResponse(resp *http.Response) *Response {
	if resp == nil {
		return &Response{
			RawResponse: resp,
		}
	}
	return &Response{
		StatusCode:  resp.StatusCode,
		Header:      resp.Header,
		RawResponse: resp,
	}
}
