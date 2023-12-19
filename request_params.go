package jac

import (
	"time"
)

// RequestParams is a structure for performing different requests
type RequestParams struct {
	method      string
	Endpoint    string
	Body        []byte
	Destination any
	Query       map[string]string
	Header      map[string]string
	Timeout     time.Duration
}

func (rp RequestParams) addMethod(method string) RequestParams {
	rp.method = method
	return rp
}
