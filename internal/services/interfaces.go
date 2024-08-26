package services

import (
	"context"
	"time"

	"github.com/v1adhope/auth-service/internal/models"
)

type Alerter interface {
	Do(email, msg string) error
}

type AuthRepo interface {
	StoreToken(ctx context.Context, id, token string, now time.Time) error
	GetToken(ctx context.Context, id string) (string, error)
	DestroyToken(ctx context.Context, id string) error
}

type Hasher interface {
	Do(target string) (string, error)
	Check(hashedTarget, target string) error
}

type TokenManager interface {
	GeneratePair(ip string, userId string) (models.TokenPair, error)
	ExtractRefreshPayload(token string) (string, error)
	ExtractAccessPayload(token string) (id, ip, userId string, err error)
}

type Validater interface {
	ValidateGuid(target string) error
}
