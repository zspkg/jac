package examples

import (
	"encoding/json"

	"github.com/zspkg/jac"
)

// FooServiceConnector is your custom service connector
type FooServiceConnector struct {
	jac.Jac
	createFooEndpoint string
}

// Foo is the service object
type Foo struct {
	bar int
}

// FooCreateResponse is the service response when creating Foo
type FooCreateResponse struct {
	id  int
	bar int
}

// NewFooServiceConnector is your custom connector to FooService where
// you can define your own data transformations and operations and
// then use jac.Jac's methods to easily send POST/GET/DELETE methods
func NewFooServiceConnector(baseEndpoint string) *FooServiceConnector {
	return &FooServiceConnector{
		jac.NewJac(baseEndpoint),
		"foo/create",
	}
}

// CreateFoo is an example of simple connector function
// which uses jac.Jac to create new Foo instance via connector
func (c *FooServiceConnector) CreateFoo(foo Foo) (*FooCreateResponse, error) {
	// rawing our Foo model
	rawFoo, err := json.Marshal(foo)
	if err != nil {
		// your custom error handling
	}

	// creating response variable
	var response FooCreateResponse

	// sending POST request to our service via connector
	// to create new Foo instance
	apiErrs, err := c.Post(
		jac.RequestParams{
			Endpoint: c.createFooEndpoint,
			Body:     rawFoo,
		},
		&response,
	)
	if err != nil {
		// your custom error handling
	}
	if len(apiErrs) != 0 {
		// you can handle API errs as you wish or just ignore them
	}

	// if err and apiErrs == nil, then you got a desired response from your service.
	// Now you can simply return it
	return &response, nil
}
