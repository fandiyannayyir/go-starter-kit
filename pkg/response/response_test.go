package response_test

import (
	"fmt"
	"math"
	"net/http"
	"testing"
	"testing/quick"

	"go-starter-kit/pkg/response"
)

func TestResponseEnvelopeStructureProperty(t *testing.T) {
	property := func(payload string) bool {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		ok := response.OK(req, payload, "ok")
		err := response.Error(req, "COMMON_INTERNAL", "failed")

		return ok.Success &&
			ok.Data == payload &&
			ok.Error == nil &&
			!err.Success &&
			err.Error != nil &&
			err.Error.Code == "COMMON_INTERNAL" &&
			err.Error.Message == "failed"
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestPaginationMetadataCalculationProperty(t *testing.T) {
	property := func(totalRaw uint16, pageRaw uint8, perPageRaw uint8) bool {
		total := int(totalRaw)
		page := int(pageRaw)
		perPage := int(perPageRaw)
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		got := response.Paginated[string](req, nil, page, perPage, total, "ok")

		expectedPage := page
		if expectedPage < 1 {
			expectedPage = 1
		}
		expectedPerPage := perPage
		if expectedPerPage <= 0 {
			expectedPerPage = 10
		}
		if expectedPerPage > 100 {
			expectedPerPage = 100
		}
		expectedTotalPages := 0
		if total > 0 {
			expectedTotalPages = int(math.Ceil(float64(total) / float64(expectedPerPage)))
		}

		return got.Meta != nil &&
			got.Meta.Page == expectedPage &&
			got.Meta.PerPage == expectedPerPage &&
			got.Meta.Total == total &&
			got.Meta.TotalPages == expectedTotalPages
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestValidationErrorCompletenessProperty(t *testing.T) {
	property := func(countRaw uint8) bool {
		count := int(countRaw % 50)
		fields := make([]response.FieldErr, 0, count)
		for i := 0; i < count; i++ {
			fields = append(fields, response.FieldErr{
				Field:   fmt.Sprintf("field_%d", i),
				Message: "required",
			})
		}
		req, _ := http.NewRequest(http.MethodPost, "/", nil)

		got := response.ValidationError(req, "COMMON_VALIDATION", "invalid", fields)

		return got.Error != nil && len(got.Error.Fields) == count
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRequestIDPropagationProperty(t *testing.T) {
	property := func(requestID string) bool {
		if requestID == "" {
			requestID = "request-id"
		}
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", requestID)

		got := response.OK(req, "payload", "ok")

		return got.RequestID == requestID
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatal(err)
	}
}
