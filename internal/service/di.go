// Package service aggregates the domain service subpackages (one folder per
// domain) into the dig container. Add a new domain by providing its constructor
// here.
package service

import (
	"go.uber.org/dig"

	"go-starter-kit/internal/service/example"
)

// Register provides every domain service's constructor to the container. dig
// injects each service's Repository (provided by internal/repository).
func Register(c *dig.Container) error {
	return c.Provide(example.NewService)
}
