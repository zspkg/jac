package jac

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

var (
	// jacDefaultConfigKey is a key in .config file corresponding
	// to the Jac configuration
	jacDefaultConfigKey = "jac"
)

// jacer is a struct implementing JACer interface
type jacer struct {
	once   comfig.Once
	getter kv.Getter
}

// JacConfig contains configurable data of a Jac
type JacConfig struct {
	URL string `fig:"url,required"`
}

// NewJACer returns an instance of JACer structure that configures Jac
// based on a provided config from kv.Getter
func NewJACer(getter kv.Getter) JACer {
	return &jacer{getter: getter}
}

// GetJacConfig returns Jac configuration info based on a provided config from kv.Getter
//   - configKey is a key in .config file corresponding to the Jac configuration.
//     If nil, then default key is used
func (c *jacer) GetJacConfig(configKey *string) JacConfig {
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
}

// ConfigureJac returns configured Jac based on a provided config from kv.Getter
//   - configKey is a key in .config file corresponding to the Jac configuration.
//     If nil, then default key is used
func (c *jacer) ConfigureJac(configKey *string) Jac {
	cfg := c.GetJacConfig(configKey)
	return NewJac(cfg.URL)
}
