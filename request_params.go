package jac

import (
	"net/http"
)

// RequestParams is a structure for performing different requests
type RequestParams struct {
	method   string
	Endpoint string
	Body     []byte
	Query    map[string]string
	Header   map[string]string
}

func (rp RequestParams) addMethod(method string) RequestParams {
	rp.method = method
	return rp
}

func (rp RequestParams) addRequestQuery(r *http.Request) *http.Request {
	if rp.Query != nil {
		q := r.URL.Query()

		for key, value := range rp.Query {
			q.Add(key, value)
		}

		r.URL.RawQuery = q.Encode()
	}

	return r
}

func (rp RequestParams) addRequestHeaders(r *http.Request) *http.Request {
	if rp.Header != nil {
		for key, value := range rp.Header {
			r.Header.Set(key, value)
		}
	}

	return r
}
