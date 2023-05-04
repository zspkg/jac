package jac

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"io"
	"net/http"
)

// jac is a structure that implements Jac interface
type jac struct {
	BaseUrl string
	JWT     *string
}

// NewJac returns new jac instance that implements Jac interface
func NewJac(baseUrl string, jwt *string) Jac {
	return &jac{baseUrl, jwt}
}

// readResponseBody reads response body into destination and returns
// respErrsPayload in case of API errors with status code higher than 400
// or err in case of some other problem happened
func (c *jac) readResponseBody(response *http.Response, destination any) (
	respErrsPayload *jsonapi.ErrorsPayload,
	err error,
) {
	// closing response body
	defer func(Body io.ReadCloser) {
		if tempErr := Body.Close(); tempErr != nil {
			err = tempErr
		}
	}(response.Body)

	// parsing response
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	// if status code is equal or higher than BadRequest
	// we are unmarshalling into errors payload
	if response.StatusCode >= http.StatusBadRequest {
		var errsPayload jsonapi.ErrorsPayload
		err = json.Unmarshal(raw, &errsPayload)
		return &errsPayload, err
	}

	// if destination is nil, we do not read response body
	if destination == nil {
		return
	}

	err = json.Unmarshal(raw, &destination)
	return nil, err
}
