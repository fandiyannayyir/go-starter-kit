package di

import (
	"go.uber.org/dig"

	"go-starter-kit/internal/handler"
	"go-starter-kit/internal/repository"
	"go-starter-kit/internal/service"
)

// RegisterDependencies wires the whole graph into the container: infra providers
// first, then each layer's constructors. dig resolves the rest (repo -> service
// -> handler) on demand when cmd invokes.
func RegisterDependencies(container *dig.Container) error {
	if err := RegisterDependency(container); err != nil {
		return err
	}
	if err := repository.Register(container); err != nil {
		return err
	}
	if err := service.Register(container); err != nil {
		return err
	}
	if err := handler.Register(container); err != nil {
		return err
	}
	return nil
}
