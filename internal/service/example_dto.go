package service

import (
	"time"

	"go-starter-kit/internal/domain/example"
)

// DTOs: the input/output contract of ExampleService (mazhab A — the service
// owns its request/response shapes; handlers just bind, validate and pass).

type CreateExampleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateExampleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExampleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toExampleResponse(e *example.Example) ExampleResponse {
	return ExampleResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
