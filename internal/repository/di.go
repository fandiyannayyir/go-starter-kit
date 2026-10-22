package repository

import "go.uber.org/dig"

// Register provides this layer's constructors to the container. dig resolves the
// *gorm.DB they need from the infra providers in internal/di.
func Register(c *dig.Container) error {
	return c.Provide(NewExampleRepository)
}
