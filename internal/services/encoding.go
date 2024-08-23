package services

import (
	"encoding/base64"

	"github.com/v1adhope/auth-service/internal/models"
)

func EncodeBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func DecodeBase64(encodetext string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(encodetext)
	if err != nil {
		return "", models.ErrNotValidTokens
	}

	return string(text), nil
}
