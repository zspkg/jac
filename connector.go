package jac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// jac is a structure that implements Jac interface
type jac struct {
	BaseUrl string
	JWT     *string
	client  *http.Client
}

// NewJac returns new jac instance that implements Jac interface
func NewJac(baseUrl string, jwt *string) Jac {
	return &jac{baseUrl, jwt, http.DefaultClient}
}

func (c *jac) Get(endpoint string, destination any) ([]*jsonapi.ErrorObject, error) {
	return c.perform(http.MethodGet, endpoint, nil, destination)
}

func (c *jac) Post(endpoint string, data []byte, destination any) ([]*jsonapi.ErrorObject, error) {
	return c.perform(http.MethodPost, endpoint, data, destination)
}

func (c *jac) Patch(endpoint string, data []byte, destination any) ([]*jsonapi.ErrorObject, error) {
	return c.perform(http.MethodPatch, endpoint, data, destination)
}

func (c *jac) Delete(endpoint string) ([]*jsonapi.ErrorObject, error) {
	return c.perform(http.MethodDelete, endpoint, nil, nil)
}

func (c *jac) Exists(endpoint string) (bool, error) {
	jsonErrs, err := c.Get(endpoint, nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to validate if object exists")
	}

	for _, jsonErr := range jsonErrs {
		if jsonErr != nil {
			if jsonErr.Status == fmt.Sprintf("%v", http.StatusNotFound) {
				return false, err
			}
		}

		return false, errors.New(
			fmt.Sprintf("unexpected error with status code %s and detail %s", jsonErr.Status, jsonErr.Detail),
		)
	}

	return true, err
}

func (c *jac) NotExists(endpoint string) (bool, error) {
	exists, err := c.Exists(endpoint)
	return exists == false, err
}

// perform performs a request based on given parameters
func (c *jac) perform(method, endpoint string, data []byte, destination any) ([]*jsonapi.ErrorObject, error) {
	resolvedEndpoint, err := c.resolveEndpoint(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve endpoint")
	}

	response, err := c.do(method, resolvedEndpoint, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	errsPayload, err := c.readResponseBody(response, destination)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if errsPayload != nil {
		return errsPayload.Errors, err
	}

	return nil, nil
}

// resolveEndpoint forms url by adding endpoint to base url.
// It ignores possible errors
func (c *jac) resolveEndpoint(endpoint string) (string, error) {
	joinedPath, err := url.JoinPath(c.BaseUrl, endpoint)
	if err != nil {
		return "", errors.Wrap(err, "failed to join path", logan.F{
			"base":     c.BaseUrl,
			"endpoint": endpoint,
		})
	}

	decodedURL, err := url.QueryUnescape(joinedPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode escape URL characters", logan.F{
			"joined_path": joinedPath,
		})
	}

	return decodedURL, nil
}

// do sends specified request to specified endpoint based on received method and data
func (c *jac) do(method, endpoint string, data []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a request")
	}

	request = c.setAuthorization(request)

	return c.client.Do(request)
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

// setAuthorization sets JWT to the Authorization header.
// If no JWT token were provided, function simply returns unmodified request
func (c *jac) setAuthorization(r *http.Request) *http.Request {
	if c.JWT == nil {
		return r
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *c.JWT))
	return r
}
