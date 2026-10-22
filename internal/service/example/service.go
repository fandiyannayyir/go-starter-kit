// Package example holds the example domain's application service: its
// repository port, DTOs and business logic. Copy this package as the starting
// point for a new domain's service.
package example

import (
	"context"

	"github.com/google/uuid"

	domain "go-starter-kit/internal/domain/example"
)

// Service holds the example business logic. Operations spanning more than one
// write should wrap their repo calls in shared.TxManager.WithTransaction so
// GetDB in the repository picks up the transaction from the context.
type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Response, error)
	Get(ctx context.Context, id string) (*Response, error)
	List(ctx context.Context, offset, limit int) ([]Response, int, error)
	Update(ctx context.Context, id string, req UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*Response, error) {
	if req.Name == "" {
		return nil, domain.ErrNameRequired
	}
	e := &domain.Example{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	resp := toResponse(e)
	return &resp, nil
}

func (s *service) Get(ctx context.Context, id string) (*Response, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toResponse(e)
	return &resp, nil
}

func (s *service) List(ctx context.Context, offset, limit int) ([]Response, int, error) {
	items, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]Response, len(items))
	for i := range items {
		resp[i] = toResponse(&items[i])
	}
	return resp, int(total), nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateRequest) (*Response, error) {
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
	resp := toResponse(e)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
