package apperror

import (
	"fmt"
	"net/http"
	"strings"
)

type AppError interface {
	error
	Code() string
	HTTPStatus() int
	Domain() string
	Unwrap() error
}

type domainError struct {
	domain     string
	code       string
	message    string
	httpStatus int
	err        error
}

func New(domain, suffix, message string, httpStatus int) AppError {
	return &domainError{
		domain:     strings.ToLower(domain),
		code:       buildCode(domain, suffix),
		message:    message,
		httpStatus: normalizeHTTPStatus(httpStatus),
	}
}

func Wrap(domain, suffix, message string, httpStatus int, err error) AppError {
	return &domainError{
		domain:     strings.ToLower(domain),
		code:       buildCode(domain, suffix),
		message:    message,
		httpStatus: normalizeHTTPStatus(httpStatus),
		err:        err,
	}
}

func (e *domainError) Error() string {
	if e.err == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %v", e.message, e.err)
}

func (e *domainError) Code() string {
	return e.code
}

func (e *domainError) HTTPStatus() int {
	return e.httpStatus
}

func (e *domainError) Domain() string {
	return e.domain
}

func (e *domainError) Unwrap() error {
	return e.err
}

func buildCode(domain, suffix string) string {
	domain = strings.TrimSpace(domain)
	suffix = strings.TrimSpace(suffix)
	return strings.ToUpper(domain + "_" + suffix)
}

func normalizeHTTPStatus(status int) int {
	if status < 100 || status > 599 {
		return http.StatusInternalServerError
	}
	return status
}
