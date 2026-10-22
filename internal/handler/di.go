package handler

import "go.uber.org/dig"

// Register provides this layer's constructors to the container. dig injects each
// handler's service (provided by internal/service) by constructor.
func Register(c *dig.Container) error {
	if err := c.Provide(NewHealthHandler); err != nil {
		return err
	}
	return c.Provide(NewExampleHandler)
}
