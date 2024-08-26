package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v1adhope/auth-service/internal/models"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/validator"
)

func setUp() *validator.Validator {
	return validator.New()
}

func TestValidateGuid(t *testing.T) {
	v := setUp()

	tcs := []struct {
		key   string
		input string
	}{
		{
			key:   "Case 1",
			input: "78ea5084-70a6-4568-8a4f-f18f3682205a",
		},
		{
			key:   "Case 2",
			input: "969c1128-d201-442a-92d7-5e7b635d11d7",
		},
		{
			key:   "Case 3",
			input: "138978a3-08bb-482b-b98a-ab2d95268dcf",
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			sut := v.ValidateGuid(tc.input)

			assert.NoError(t, sut, tc.key)
		})
	}
}

func TestValidateGuidNegative(t *testing.T) {
	v := setUp()

	tcs := []struct {
		key   string
		input string
	}{
		{
			key:   "Case 1",
			input: "78ea5084-4568-8a4f-f18f3682205a",
		},
		{
			key:   "Case 2",
			input: "969c1128-d201-442a-92d7-5e7b635d11d77",
		},
		{
			key:   "Case 3",
			input: "11",
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			sut := v.ValidateGuid(tc.input)

			assert.ErrorIs(t, sut, models.ErrNotValidGuid, tc.key)
		})
	}
}
