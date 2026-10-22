package example

import (
	"time"

	domain "go-starter-kit/internal/domain/example"
)

// DTOs: the input/output contract of Service (mazhab A — the service owns its
// request/response shapes; handlers just bind, validate and pass).

type CreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toResponse(e *domain.Example) Response {
	return Response{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
