package hash

import "golang.org/x/crypto/bcrypt"

type Hash struct{}

func New() *Hash {
	return &Hash{}
}

func (h *Hash) Do(target string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(target), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}
