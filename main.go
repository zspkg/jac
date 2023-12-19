package jac

import "github.com/google/jsonapi"

// Jac is the interface that connector should implement
type Jac interface {
	// Get sends GET request and reads response body into destination.
	// Returns a slice of API error objects according to JSON API or
	// error if some happened during the operation.
	Get(params RequestParams, destination any) ([]*jsonapi.ErrorObject, error)
	// Post sends POST request with provided data as a request body
	// and reads response body if some data is expected to return.
	// Returns a slice of API error objects according to JSON API or
	// error if some happened during the operation.
	Post(params RequestParams, destination any) ([]*jsonapi.ErrorObject, error)
	// Patch sends PATCH request with provided data as a request body
	// and reads response body if some data is expected to return.
	// Returns a slice of API error objects according to JSON API or
	// error if some happened during the operation.
	Patch(params RequestParams, destination any) ([]*jsonapi.ErrorObject, error)
	// Delete sends DELETE request.
	// Returns a slice of API error objects according to JSON API or
	// error if some happened during the operation.
	Delete(params RequestParams) ([]*jsonapi.ErrorObject, error)
	// Exists checks if object exists by provided endpoint.
	// Returns error if non-2xx status differs from 404 or
	// something happened during the operation.
	Exists(params RequestParams) (bool, error)
	// NotExists checks if object is not exist by provided endpoint.
	// Returns error if non-2xx status differs from 404 or
	// something happened during the operation.
	NotExists(params RequestParams) (bool, error)
}

// JACer is the interface that connector configurator should implement
type JACer interface {
	//GetJacConfig returns Jac configuration info based on config information
	// that can be found by specified key. It is useful when you need
	// to configure multiple connectors to different services.
	// If jacConfigKey is nil, default config key ``jac`` is used.
	GetJacConfig(jacConfigKey *string) JacConfig
	// ConfigureJac returns configured Jac based on config information
	// that can be found by specified key. It is useful when you need
	// to configure multiple connectors to different services.
	// If jacConfigKey is nil, default config key ``jac`` is used.
	ConfigureJac(jacConfigKey *string) Jac
}
