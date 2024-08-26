package hash

import (
	"fmt"

	"github.com/v1adhope/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func New() *Hash {
	return &Hash{}
}

func (h *Hash) Do(target string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(target), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash: hash: Do: GenerateFromPassword: %w", models.ErrNotValidTokens)
	}

	return string(hashBytes), nil
}

func (h *Hash) Check(hashedTarget, target string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedTarget), []byte(target))
}
