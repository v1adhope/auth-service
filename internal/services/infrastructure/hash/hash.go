package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func New() *Hash {
	return &Hash{}
}

func (h *Hash) Do(target string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(target), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash: hash: Do: GenerateFromPassword: %w", err)
	}

	return string(hashBytes), nil
}
