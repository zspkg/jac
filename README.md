# :jack_o_lantern: JAC — JSON API service connector building base
[![Go Reference](https://pkg.go.dev/badge/github.com/zspkg/jac#section-readme.svg)](https://pkg.go.dev/github.com/zspkg/jac#section-readme)
[![codecov](https://codecov.io/github/zspkg/jac/branch/main/graph/badge.svg?token=JO5Qd0Zw20)](https://codecov.io/github/zspkg/jac)
[![Go Report Card](https://goreportcard.com/badge/github.com/zspkg/jac)](https://goreportcard.com/report/github.com/zspkg/jac)

JAC is a building base for custom JSON API service connectors. 

## Usage example

```go
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
func NewFooServiceConnector(baseEndpoint string, jwt *string) *FooServiceConnector {
	return &FooServiceConnector{
		jac.NewJac(baseEndpoint, jwt),
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
	apiErrs, err := c.Post(c.createFooEndpoint, rawFoo, &response)
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
```

Note that you can configure `Jac` directly from config using `JACer`. It uses `Getter` which is responsible for retrieving info from config files and must implement next interface:
```go
type Getter interface {
    GetStringMap(key string) (map[string]interface{}, error)
}
```

`JACer` has next methods to configure `Jac` connector or to simply retrieve `Jac` configuration:

```go
// JacConfig contains configurable data of a Jac
type JacConfig struct {
	URL string  `fig:"url,required"`
	JWT *string `fig:"jwt"`
}

// NewJACer returns an instance of JACer structure that configures Jac
func NewJACer(getter kv.Getter) JACer {
	return &jacer{getter: getter}
}

// GetJacConfig returns Jac configuration info based on a provided config from kv.Getter
func (c *jacer) GetJacConfig(configKey *string) JacConfig {
	return c.once.Do(func() interface{} {
		if configKey == nil {
			configKey = &jacDefaultConfigKey
		}

		var (
			config = JacConfig{}
			raw    = kv.MustGetStringMap(c.getter, *configKey)
		)

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out jac"))
		}

		return config
	}).(JacConfig)
}

// ConfigureJac returns configured Jac based on a provided config from kv.Getter
func (c *jacer) ConfigureJac(configKey *string) Jac {
	cfg := c.GetJacConfig(configKey)
	return NewJac(cfg.URL, cfg.JWT)
}
```
