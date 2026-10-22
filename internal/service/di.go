package service

import "go.uber.org/dig"

// Register provides this layer's constructors to the container. dig injects each
// service's repository (provided by internal/repository) by its interface type.
func Register(c *dig.Container) error {
	return c.Provide(NewExampleService)
}
