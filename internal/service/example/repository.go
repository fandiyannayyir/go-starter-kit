package example

import (
	"context"

	domain "go-starter-kit/internal/domain/example"
)

// Repository is the persistence port this service depends on. It is declared on
// the consumer side (Go idiom) so the service can be unit-tested with a fake and
// the backend can be swapped without touching the service. The GORM
// implementation lives in internal/repository.
type Repository interface {
	Create(ctx context.Context, e *domain.Example) error
	FindByID(ctx context.Context, id string) (*domain.Example, error)
	List(ctx context.Context, offset, limit int) ([]domain.Example, int64, error)
	Update(ctx context.Context, e *domain.Example) error
	Delete(ctx context.Context, id string) error
}
