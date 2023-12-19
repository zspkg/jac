package jac

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	jacTestConfigKey1 = "test-config-1.yaml"
	jacTestConfigKey2 = "test-config-2.yaml"
)

func TestJacer_GetJacConfig(t *testing.T) {
	t.Run("using test-config-1.yaml", func(t *testing.T) {
		myJacer := NewJACer(kv.NewViperFile(jacTestConfigKey1))
		jacCfg := myJacer.GetJacConfig(nil)

		assert.Equal(t, JacConfig{
			URL: "http://localhost:8000",
		}, jacCfg)
	})

	t.Run("using test-config-2.yaml", func(t *testing.T) {
		myJacer := NewJACer(kv.NewViperFile(jacTestConfigKey2))
		jacCfgKey := "my-connector-name"
		jacCfg := myJacer.GetJacConfig(&jacCfgKey)

		assert.Equal(t, JacConfig{
			URL: "http://localhost:8001",
		}, jacCfg)
	})

	t.Run("using non-existent config: expect panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()

		myWrongJacer := NewJACer(kv.NewViperFile("some-non-existent-config.yaml"))
		_ = myWrongJacer.GetJacConfig(nil)
	})
}
