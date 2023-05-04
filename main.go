package jac

// Jac is the interface that connector should implement
type Jac interface {
}

// JACer is the interface that connector configurator should implement
type JACer interface {
	// ConfigureJac returns configured Jac
	ConfigureJac() Jac
}
