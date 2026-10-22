package service

import (
	"context"
	"errors"
	"testing"

	"go-starter-kit/internal/domain/example"
)

// fakeExampleRepo is an in-memory ExampleRepository. Its existence is the point
// of declaring the repository interface here: the service is tested with no DB.
type fakeExampleRepo struct {
	store map[string]*example.Example
}

func newFakeRepo() *fakeExampleRepo {
	return &fakeExampleRepo{store: map[string]*example.Example{}}
}

func (f *fakeExampleRepo) Create(_ context.Context, e *example.Example) error {
	f.store[e.ID] = e
	return nil
}

func (f *fakeExampleRepo) FindByID(_ context.Context, id string) (*example.Example, error) {
	e, ok := f.store[id]
	if !ok {
		return nil, example.ErrNotFound
	}
	return e, nil
}

func (f *fakeExampleRepo) List(_ context.Context, _, _ int) ([]example.Example, int64, error) {
	out := make([]example.Example, 0, len(f.store))
	for _, e := range f.store {
		out = append(out, *e)
	}
	return out, int64(len(out)), nil
}

func (f *fakeExampleRepo) Update(_ context.Context, e *example.Example) error {
	if _, ok := f.store[e.ID]; !ok {
		return example.ErrNotFound
	}
	f.store[e.ID] = e
	return nil
}

func (f *fakeExampleRepo) Delete(_ context.Context, id string) error {
	if _, ok := f.store[id]; !ok {
		return example.ErrNotFound
	}
	delete(f.store, id)
	return nil
}

func TestCreate_EmptyNameRejected(t *testing.T) {
	svc := NewExampleService(newFakeRepo())
	_, err := svc.Create(context.Background(), CreateExampleRequest{Name: ""})
	if !errors.Is(err, example.ErrNameRequired) {
		t.Fatalf("want ErrNameRequired, got %v", err)
	}
}

func TestCreate_ThenGet(t *testing.T) {
	svc := NewExampleService(newFakeRepo())
	created, err := svc.Create(context.Background(), CreateExampleRequest{Name: "first"})
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
	svc := NewExampleService(newFakeRepo())
	_, err := svc.Get(context.Background(), "nope")
	if !errors.Is(err, example.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
