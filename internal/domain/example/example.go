// Package example is the example domain: a pure entity plus its domain errors,
// with no GORM/DB dependency. Copy this package as the starting point for a new
// product domain.
package example

import (
	"net/http"
	"time"

	"go-starter-kit/pkg/apperror"
)

// Domain is the error-code namespace for this domain (e.g. EXAMPLE_NOT_FOUND).
const Domain = "example"

// Example is the domain entity.
type Example struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrNotFound     = apperror.New(Domain, "not_found", "example not found", http.StatusNotFound)
	ErrNameRequired = apperror.New(Domain, "name_required", "name is required", http.StatusBadRequest)
)
