package example

import (
	"context"
	"errors"
	"testing"

	domain "go-starter-kit/internal/domain/example"
)

// fakeRepo is an in-memory Repository. Its existence is the point of declaring
// the repository interface here: the service is tested with no DB.
type fakeRepo struct {
	store map[string]*domain.Example
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{store: map[string]*domain.Example{}}
}

func (f *fakeRepo) Create(_ context.Context, e *domain.Example) error {
	f.store[e.ID] = e
	return nil
}

func (f *fakeRepo) FindByID(_ context.Context, id string) (*domain.Example, error) {
	e, ok := f.store[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return e, nil
}

func (f *fakeRepo) List(_ context.Context, _, _ int) ([]domain.Example, int64, error) {
	out := make([]domain.Example, 0, len(f.store))
	for _, e := range f.store {
		out = append(out, *e)
	}
	return out, int64(len(out)), nil
}

func (f *fakeRepo) Update(_ context.Context, e *domain.Example) error {
	if _, ok := f.store[e.ID]; !ok {
		return domain.ErrNotFound
	}
	f.store[e.ID] = e
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id string) error {
	if _, ok := f.store[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.store, id)
	return nil
}

func TestCreate_EmptyNameRejected(t *testing.T) {
	svc := NewService(newFakeRepo())
	_, err := svc.Create(context.Background(), CreateRequest{Name: ""})
	if !errors.Is(err, domain.ErrNameRequired) {
		t.Fatalf("want ErrNameRequired, got %v", err)
	}
}

func TestCreate_ThenGet(t *testing.T) {
	svc := NewService(newFakeRepo())
	created, err := svc.Create(context.Background(), CreateRequest{Name: "first"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected a generated id")
	}
	got, err := svc.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "first" {
		t.Fatalf("want name 'first', got %q", got.Name)
	}
}

func TestGet_MissingReturnsNotFound(t *testing.T) {
	svc := NewService(newFakeRepo())
	_, err := svc.Get(context.Background(), "nope")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
