package di

import "go.uber.org/dig"

func RegisterDependencies(container *dig.Container) error {
	if err := RegisterDependency(container); err != nil {
		return err
	}
	return nil
}
