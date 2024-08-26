package app

import (
	"context"

	"github.com/v1adhope/auth-service/internal/services"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/alert"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/hash"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/repositories"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/tokens"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/validator"
	httpv1 "github.com/v1adhope/auth-service/internal/transports/http/v1"
	"github.com/v1adhope/auth-service/pkg/httpserver"
	"github.com/v1adhope/auth-service/pkg/logger"
	"github.com/v1adhope/auth-service/pkg/postgresql"
)

func Run(ctx context.Context, cfg Config) error {
	validator := validator.New()

	tokenManager := tokens.New(
		tokens.WithAccessKey(cfg.Tokens.AccessKey),
		tokens.WithAccessTtl(cfg.Tokens.AceessTtl),
		tokens.WithRefreshKey(cfg.Tokens.RefreshKey),
		tokens.WithIssuer(cfg.Tokens.Issuer),
	)

	hash := hash.New()

	postgres, err := postgresql.Build(
		ctx,
		postgresql.WithConnStr(cfg.Postgres.ConnStr),
	)
	if err != nil {
		return err
	}
	defer postgres.Close()

	repos := repositories.New(postgres)

	alert := alert.New()

	services := services.New(
		validator,
		tokenManager,
		hash,
		repos,
		alert,
	)

	log := logger.New(
		logger.WithLevel(cfg.Logger.Level),
	)

	handler := httpv1.New(services, log).Handler(
		httpv1.WithAllowOrigins(cfg.Server.AllowOrigins),
		httpv1.WithAllowMethods(cfg.Server.AllowMethods),
		httpv1.WithAllowHeaders(cfg.Server.AllowHeaders),
		httpv1.WithMode(cfg.Server.Mode),
	)

	s := httpserver.New(
		handler,
		httpserver.WithSocket(cfg.Server.Socket),
		httpserver.WithShutdownTimeout(cfg.Server.ShutdownTimeout),
		httpserver.WithWriteTimeout(cfg.Server.WriteTimeout),
		httpserver.WithReadTimeout(cfg.Server.ReadTimeout),
	)

	s.Run()

	return nil
}
