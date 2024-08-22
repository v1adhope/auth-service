package services

import (
	"context"

	"github.com/v1adhope/auth-service/internal/models"
)

type Allerter interface {
	Do(email, msg string) error
}

type AuthRepo interface {
	Store(ctx context.Context, token string) error
	Check(ctx context.Context, token string) error
	Destroy(ctx context.Context, token string) error
}

type Hasher interface {
	Do(target string) (string, error)
}

type TokenManager interface {
	GeneratePair(id string, ip string) (models.TokenPair, error)
	RefreshPair(tp models.TokenPair, ip string) (newTp models.TokenPair, isIpChanged bool, err error)
}

type Validater interface {
	ValidateGuid(target string) error
}
