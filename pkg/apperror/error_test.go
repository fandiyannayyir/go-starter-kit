package apperror_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"testing/quick"

	"go-starter-kit/pkg/apperror"
)

func TestDomainErrorCodePrefixProperty(t *testing.T) {
	property := func(domainInput, suffixInput string) bool {
		domain := tokenOrDefault(domainInput, "auth")
		suffix := tokenOrDefault(suffixInput, "invalid_credentials")

		err := apperror.New(domain, suffix, "message", http.StatusBadRequest)

		return err.Code() == strings.ToUpper(domain+"_"+suffix)
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestErrorWrappingPreservesChainProperty(t *testing.T) {
	sentinel := errors.New("database unavailable")

	property := func(domainInput, suffixInput string) bool {
		domain := tokenOrDefault(domainInput, "auth")
		suffix := tokenOrDefault(suffixInput, "internal")

		err := apperror.Wrap(domain, suffix, "wrapped", http.StatusInternalServerError, sentinel)

		var appErr apperror.AppError
		return errors.Is(err, sentinel) && errors.As(err, &appErr) && appErr.Code() == err.Code()
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func tokenOrDefault(input, fallback string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(input) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return fallback
	}
	return b.String()
}
