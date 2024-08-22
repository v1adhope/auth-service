package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/v1adhope/auth-service/internal/models"
	"github.com/v1adhope/auth-service/pkg/postgresql"
)

type Auth struct {
	*postgresql.Postgres
}

func (r *Auth) Store(ctx context.Context, token string) error {
	sql, args, err := r.Builder.Insert("auth_whitelist").
		SetMap(squirrel.Eq{
			"token": token,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("repositories: auth: Store: ToSql: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("repositories: auth: Store: Exec: %w", err)
	}

	return nil
}

func (r *Auth) Check(ctx context.Context, token string) error {
	sql, args, err := r.Builder.Select("1").
		Prefix("select exists (").
		From("auth_whitelist").
		Where(squirrel.Eq{
			"token": token,
		}).
		Suffix(")").
		ToSql()
	if err != nil {
		return fmt.Errorf("repositories: auth: Check: ToSql: %w", err)
	}

	ok := true

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&ok); err != nil {
		return fmt.Errorf("repositories: auth: Check: Scan: %w", err)
	}

	if !ok {
		return models.ErrNotValidTokens
	}

	return nil
}

func (r *Auth) Destroy(ctx context.Context, token string) error {
	sql, args, err := r.Builder.Delete("auth_whitelist").
		Where(squirrel.Eq{
			"token": token,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repositories: auth: Destroy: ToSql: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("repositories: auth: Destroy: Exec: %w", err)
	}

	return nil
}
