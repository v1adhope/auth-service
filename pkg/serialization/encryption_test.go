package serialization_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v1adhope/auth-service/pkg/serialization"
)

type testSerializationByGcmInput struct {
	key  []byte
	text []byte
}

func TestSerializationByGcm(t *testing.T) {
	tcs := []struct {
		key   string
		input testSerializationByGcmInput
	}{
		{
			key: "Case 1",
			input: testSerializationByGcmInput{
				[]byte("vhZ35oAnPtqyu2dN"),
				[]byte("310b5d1e-eff0-4fb0-97d8-6734682320b3"),
			},
		},
		{
			key: "Case 2",
			input: testSerializationByGcmInput{
				[]byte("WSDrDbwTMFwM5QUndxhqVEqqpgzcuR9g"),
				[]byte("4b7d46a6-ab5a-4a9b-8f48-acc446be8d04"),
			},
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			ciphertext, err := serialization.EncryptByGcm(tc.input.text, tc.input.key)

			assert.NoError(t, err, tc.key)

			sut, err := serialization.DecryptByGcm(ciphertext, tc.input.key)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, tc.input.text, sut, tc.key)
		})
	}
}

type TestSerializationByGcmNegativeInput struct {
	encryptkey []byte
	decryptkey []byte
	text       []byte
}

func TestSerializationByGcmNegative(t *testing.T) {
	tcs := []struct {
		key   string
		input TestSerializationByGcmNegativeInput
	}{
		{
			key: "Case 1",
			input: TestSerializationByGcmNegativeInput{
				[]byte("vhZ35oAnPtqyu2dN"),
				[]byte("vhZ35oA2Ptqyu2dN"),
				[]byte("310b5d1e-eff0-4fb0-97d8-6734682320b3"),
			},
		},
		{
			key: "Case 2",
			input: TestSerializationByGcmNegativeInput{
				[]byte("WSDrDbwTMFwM5QUndxhqVEqqpgzcuR9g"),
				[]byte("WSDrDbwTmFwM5QUndxhqVEqqpgzcuR9g"),
				[]byte("4b7d46a6-ab5a-4a9b-8f48-acc446be8d04"),
			},
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			ciphertext, err := serialization.EncryptByGcm(tc.input.text, tc.input.encryptkey)

			assert.NoError(t, err, tc.key)

			_, sut := serialization.DecryptByGcm(ciphertext, tc.input.decryptkey)

			assert.Error(t, sut, tc.key)
		})
	}
}
