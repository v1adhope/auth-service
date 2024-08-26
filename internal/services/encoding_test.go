package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v1adhope/auth-service/internal/services"
)

func TestEncodeBase64(t *testing.T) {
	tcs := []struct {
		key      string
		input    string
		expected string
	}{
		{
			key:      "Case 1",
			input:    "caba6bb4-f3e4-4eea-9d80-6376a6e5ad84",
			expected: "Y2FiYTZiYjQtZjNlNC00ZWVhLTlkODAtNjM3NmE2ZTVhZDg0",
		},
		{
			key:      "Case 2",
			input:    "21d53ead-aa0b-4b39-b4d0-73b5f201b357",
			expected: "MjFkNTNlYWQtYWEwYi00YjM5LWI0ZDAtNzNiNWYyMDFiMzU3",
		},
		{
			key:      "Case 3",
			input:    "3422a8b6-2ec2-48c9-836e-066cd4083232",
			expected: "MzQyMmE4YjYtMmVjMi00OGM5LTgzNmUtMDY2Y2Q0MDgzMjMy",
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			sut := services.EncodeBase64(tc.input)

			assert.Equal(t, tc.expected, sut, tc.key)
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	tcs := []struct {
		key      string
		input    string
		expected string
	}{
		{
			key:      "Case 1",
			input:    "Y2FiYTZiYjQtZjNlNC00ZWVhLTlkODAtNjM3NmE2ZTVhZDg0",
			expected: "caba6bb4-f3e4-4eea-9d80-6376a6e5ad84",
		},
		{
			key:      "Case 2",
			input:    "MjFkNTNlYWQtYWEwYi00YjM5LWI0ZDAtNzNiNWYyMDFiMzU3",
			expected: "21d53ead-aa0b-4b39-b4d0-73b5f201b357",
		},
		{
			key:      "Case 3",
			input:    "MzQyMmE4YjYtMmVjMi00OGM5LTgzNmUtMDY2Y2Q0MDgzMjMy",
			expected: "3422a8b6-2ec2-48c9-836e-066cd4083232",
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			sut, err := services.DecodeBase64(tc.input)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, tc.expected, sut, tc.key)
		})
	}
}
