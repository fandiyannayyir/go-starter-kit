package response

import (
	"math"
	"net/http"

	"github.com/google/uuid"
)

type Envelope[T any] struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message,omitempty"`
	Data      T          `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	Meta      *Meta      `json:"meta,omitempty"`
	RequestID string     `json:"request_id"`
}

type ErrorBody struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Fields  []FieldErr `json:"fields,omitempty"`
}

type FieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PaginatedData[T any] struct {
	Items []T `json:"items"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func OK[T any](r *http.Request, data T, message string) Envelope[T] {
	return Envelope[T]{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: getRequestID(r),
	}
}

func Created[T any](r *http.Request, data T, message string) Envelope[T] {
	return OK(r, data, message)
}

func Error(r *http.Request, code, message string) Envelope[any] {
	return Envelope[any]{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
		RequestID: getRequestID(r),
	}
}

func ValidationError(r *http.Request, code, message string, fields []FieldErr) Envelope[any] {
	return Envelope[any]{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
		RequestID: getRequestID(r),
	}
}

func Paginated[T any](r *http.Request, items []T, page, perPage, total int, message string) Envelope[PaginatedData[T]] {
	page, perPage = normalizePagination(page, perPage)
	return Envelope[PaginatedData[T]]{
		Success: true,
		Message: message,
		Data:    PaginatedData[T]{Items: items},
		Meta: &Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      max(total, 0),
			TotalPages: totalPages(total, perPage),
		},
		RequestID: getRequestID(r),
	}
}

func getRequestID(r *http.Request) string {
	if r != nil {
		if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
			return requestID
		}
	}
	return uuid.NewString()
}

func normalizePagination(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func totalPages(total, perPage int) int {
	if total <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(perPage)))
}
