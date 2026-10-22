package service

import (
	"context"

	"github.com/google/uuid"

	"go-starter-kit/internal/domain/example"
)

// ExampleService holds the example business logic. Operations spanning more
// than one write should wrap their repo calls in shared.TxManager.WithTransaction
// so GetDB in the repository picks up the transaction from the context.
type ExampleService interface {
	Create(ctx context.Context, req CreateExampleRequest) (*ExampleResponse, error)
	Get(ctx context.Context, id string) (*ExampleResponse, error)
	List(ctx context.Context, offset, limit int) ([]ExampleResponse, int, error)
	Update(ctx context.Context, id string, req UpdateExampleRequest) (*ExampleResponse, error)
	Delete(ctx context.Context, id string) error
}

type exampleService struct {
	repo ExampleRepository
}

func NewExampleService(repo ExampleRepository) ExampleService {
	return &exampleService{repo: repo}
}

func (s *exampleService) Create(ctx context.Context, req CreateExampleRequest) (*ExampleResponse, error) {
	if req.Name == "" {
		return nil, example.ErrNameRequired
	}
	e := &example.Example{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	resp := toExampleResponse(e)
	return &resp, nil
}

func (s *exampleService) Get(ctx context.Context, id string) (*ExampleResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toExampleResponse(e)
	return &resp, nil
}

func (s *exampleService) List(ctx context.Context, offset, limit int) ([]ExampleResponse, int, error) {
	items, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]ExampleResponse, len(items))
	for i := range items {
		resp[i] = toExampleResponse(&items[i])
	}
	return resp, int(total), nil
}

func (s *exampleService) Update(ctx context.Context, id string, req UpdateExampleRequest) (*ExampleResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		e.Name = req.Name
	}
	e.Description = req.Description
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}
	resp := toExampleResponse(e)
	return &resp, nil
}

func (s *exampleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
