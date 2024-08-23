package httpv1

import (
	"context"

	"github.com/v1adhope/auth-service/internal/models"
)

type AuthService interface {
	GenerateTokenPair(ctx context.Context, userId string, ip string) (models.TokenPair, error)
	RefreshTokenPair(ctx context.Context, tp models.TokenPair, ip string) (models.TokenPair, error)
}

type Logger interface {
	Info(format string, msg ...any)
	Debug(err error, format string, msg ...any)
	Error(err error, format string, msg ...any)
}
