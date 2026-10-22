package service

import (
	"context"

	"go-starter-kit/internal/domain/example"
)

// ExampleRepository is the persistence port the example service depends on. It
// is declared here, on the consumer side (Go idiom), so the service can be
// unit-tested with a fake and the backend can be swapped without touching the
// service. The GORM implementation lives in internal/repository.
type ExampleRepository interface {
	Create(ctx context.Context, e *example.Example) error
	FindByID(ctx context.Context, id string) (*example.Example, error)
	List(ctx context.Context, offset, limit int) ([]example.Example, int64, error)
	Update(ctx context.Context, e *example.Example) error
	Delete(ctx context.Context, id string) error
}
