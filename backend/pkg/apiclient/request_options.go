package apiclient

import (
	"net/http"
)

type RequestOption func(*http.Request) error

func (c *Client) applyRequestOptions(req *http.Request, requestOptions ...RequestOption) error {
	for _, option := range requestOptions {
		if err := option(req); err != nil {
			return err
		}
	}

	return nil
}

// WithHTTPHeader is a request option function that sets a given HTTP header.
func WithHTTPHeader(name, value string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set(name, value)
		return nil
	}
}
