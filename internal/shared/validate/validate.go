// Package validate wraps go-playground/validator as an echo.Validator and
// translates its errors into response.FieldErr, so the central error handler can
// render a validation failure as the standard envelope with per-field messages.
package validate

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"go-starter-kit/pkg/response"
)

// Validator implements echo.Validator.
type Validator struct {
	v *validator.Validate
}

// New builds a Validator that reports each field by its json tag name (so error
// fields match the request body, e.g. "name" not "Name").
func New() *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{v: v}
}

// Validate satisfies echo.Validator; handlers call c.Validate(&req).
func (cv *Validator) Validate(i any) error {
	return cv.v.Struct(i)
}

// FieldErrors converts a validator error into response.FieldErr entries. ok is
// false when err is not a validation error (the caller then handles it normally).
func FieldErrors(err error) ([]response.FieldErr, bool) {
	var valErrs validator.ValidationErrors
	if !errors.As(err, &valErrs) {
		return nil, false
	}
	out := make([]response.FieldErr, 0, len(valErrs))
	for _, fe := range valErrs {
		out = append(out, response.FieldErr{Field: fe.Field(), Message: message(fe)})
	}
	return out, true
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + fe.Param()
	case "max":
		return "must be at most " + fe.Param()
	default:
		return "is invalid"
	}
}
