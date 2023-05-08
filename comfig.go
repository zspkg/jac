package jac

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	// jacConfigKey is a key in .config file corresponding
	// to the Jac configuration
	jacConfigKey = "jac"
)

// jacer is a struct implementing JACer interface
type jacer struct {
	once   comfig.Once
	getter kv.Getter
}

// jacConfig contains configurable data of a Jac
type jacConfig struct {
	URL string  `fig:"url,required"`
	JWT *string `fig:"jwt"`
}

// NewJACer returns an instance of JACer structure that configures Jac
func NewJACer(getter kv.Getter) JACer {
	return &jacer{getter: getter}
}

// ConfigureJac returns configured Jac based on a provided config from kv.Getter
func (c *jacer) ConfigureJac() Jac {
	return c.once.Do(func() interface{} {
		var (
			config = jacConfig{}
			raw    = kv.MustGetStringMap(c.getter, jacConfigKey)
		)

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out jac"))
		}

		return NewJac(config.URL, config.JWT)
	}).(Jac)
}
